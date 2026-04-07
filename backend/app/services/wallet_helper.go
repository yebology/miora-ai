// wallet_helper.go contains helper methods for WalletService.
// These are internal implementation details — not part of the public service interface.
package services

import (
	"errors"
	"time"

	"miora-ai/app/dto"
	"miora-ai/app/dto/responses"
	"miora-ai/app/entities"
	"miora-ai/app/interfaces"
	"miora-ai/config"
	"miora-ai/constants"
	"miora-ai/pkg"
	"miora-ai/utils"

	"gorm.io/gorm"
)

// --- Data fetching helpers ---

// findOrCreateWallet looks up a wallet by address, creates it if not found.
func (s *WalletService) findOrCreateWallet(address, chain string) (*entities.Wallet, *pkg.AppError) {

	wallet, err := s.repo.FindByAddress(address)
	if err == nil {
		return wallet, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, pkg.ErrInternal()
	}

	wallet = &entities.Wallet{Address: address, Chain: chain}
	if err := s.repo.Create(wallet); err != nil {
		return nil, pkg.ErrInternal()
	}

	return wallet, nil

}

// getClient returns the blockchain client for the given chain.
func (s *WalletService) getClient(chain string) (interfaces.BlockchainClient, *pkg.AppError) {

	if constants.IsEVM(chain) {
		return s.evmClient, nil
	}
	if constants.IsSolana(chain) {
		return s.svmClient, nil
	}
	return nil, pkg.ErrBadReq(constants.UnsupportedChain)

}

// fetchAndSaveTransactions fetches transfers from Alchemy and persists them.
func (s *WalletService) fetchAndSaveTransactions(wallet *entities.Wallet, chain string, limit int) ([]entities.Transaction, *pkg.AppError) {

	client, appErr := s.getClient(chain)
	if appErr != nil {
		return nil, appErr
	}

	transfers, err := client.GetTransfers(wallet.Address, limit, chain)
	if err != nil {
		return nil, pkg.ErrUnexpected(502, constants.AnalysisFailed)
	}

	txEntities := make([]entities.Transaction, 0, len(transfers))
	for _, t := range transfers {
		txEntities = append(txEntities, entities.Transaction{
			WalletID:        wallet.ID,
			Hash:            t.Hash,
			Chain:           chain,
			From:            t.From,
			To:              t.To,
			Value:           t.Value,
			TokenSymbol:     t.TokenSymbol,
			ContractAddress: t.ContractAddress,
			Direction:       t.Direction,
			BlockNumber:     t.BlockNumber,
			Timestamp:       time.Unix(t.Timestamp, 0),
		})
	}

	if err := s.repo.SaveTransactions(txEntities); err != nil {
		return nil, pkg.ErrInternal()
	}

	return txEntities, nil

}

// fetchTokenData queries DexScreener for each unique token in the transactions.
func (s *WalletService) fetchTokenData(chain string, txs []entities.Transaction) map[string]dto.TokenPairData {

	dexChain := chainToDexScreenerID(chain)
	tokenData := make(map[string]dto.TokenPairData)
	seen := make(map[string]bool)

	for _, tx := range txs {
		addr := tx.ContractAddress
		if addr == "" || seen[addr] {
			continue
		}
		seen[addr] = true

		pairs, err := s.dexScreener.GetTokenPairs(dexChain, addr)
		if err == nil && len(pairs) > 0 {
			tokenData[addr] = pairs[0]
		}
	}

	return tokenData

}

// --- Price helpers ---

// getPrice fetches token price — Moralis for EVM (by block), Birdeye for Solana (by timestamp).
// block=0 and timestamp=now means current price.
func (s *WalletService) getPrice(chain, tokenAddr string, blockNumber uint64, timestamp time.Time) float64 {

	if constants.IsSolana(chain) {
		data, err := s.birdeye.GetHistoricalPrice(tokenAddr, timestamp.Unix())
		if err == nil && data.UsdPrice > 0 {
			return data.UsdPrice
		}
	} else {
		data, err := s.moralis.GetTokenPrice(chain, tokenAddr, blockNumber)
		if err == nil && data.UsdPrice > 0 {
			return data.UsdPrice
		}
	}

	return 0

}

// --- Grouping helpers ---

type tokenGroup struct {
	buys  []entities.Transaction
	sells []entities.Transaction
}

// groupByToken groups transactions by contract address into buys and sells, sorted by block.
func (s *WalletService) groupByToken(txs []entities.Transaction) map[string]*tokenGroup {

	grouped := make(map[string]*tokenGroup)

	for _, tx := range txs {
		addr := tx.ContractAddress
		if addr == "" || tx.BlockNumber == 0 {
			continue
		}

		if _, ok := grouped[addr]; !ok {
			grouped[addr] = &tokenGroup{}
		}

		if tx.Direction == "in" {
			grouped[addr].buys = insertSorted(grouped[addr].buys, tx)
		} else if tx.Direction == "out" {
			grouped[addr].sells = insertSorted(grouped[addr].sells, tx)
		}
	}

	return grouped

}

// --- Utility functions ---

func chainToDexScreenerID(chain string) string {

	cfg := constants.GetChainConfig(chain)
	if cfg != nil {
		return cfg.DexScreenerID
	}
	return chain

}

// insertSorted inserts a transaction sorted by BlockNumber ascending.
func insertSorted(txs []entities.Transaction, tx entities.Transaction) []entities.Transaction {

	i := len(txs)
	for i > 0 && txs[i-1].BlockNumber > tx.BlockNumber {
		i--
	}
	txs = append(txs, entities.Transaction{})
	copy(txs[i+1:], txs[i:])
	txs[i] = tx
	return txs

}

// --- Response builders ---

// buildTradedTokens converts trade results into response DTOs.
// Looks up token symbol from transaction data.
func buildTradedTokens(chain string, trades []tradeResult, txs []entities.Transaction) []responses.TradedToken {

	// Build symbol lookup from transactions
	symbolMap := make(map[string]string)
	for _, tx := range txs {
		if tx.ContractAddress != "" && tx.TokenSymbol != "" {
			symbolMap[tx.ContractAddress] = tx.TokenSymbol
		}
	}

	tokens := make([]responses.TradedToken, 0, len(trades))
	for _, t := range trades {
		status := "unrealized"
		var exitTime *time.Time

		if !t.ExitTime.IsZero() {
			status = "realized"
			et := t.ExitTime
			exitTime = &et
		}

		tokens = append(tokens, responses.TradedToken{
			ContractAddress: t.TokenAddress,
			Symbol:          symbolMap[t.TokenAddress],
			Chain:           chain,
			PnlPercent:      utils.Round2(t.PnlPercent),
			BuyPrice:        t.BuyPrice,
			ExitPrice:       t.ExitPrice,
			BuyTime:         t.BuyTime,
			ExitTime:        exitTime,
			Status:          status,
		})
	}

	return tokens

}

// --- Condition generators ---

// buildConditions generates suggested follow conditions based on scoring data.
// Only generated for "conditional_follow" recommendations.
// Conditions are based on the wallet's weak scoring areas.
func buildConditions(
	tokenData map[string]dto.TokenPairData,
	riskExposure, entryTiming, tokenQuality float64,
	scoringCfg config.ScoringConfig,
) []responses.Condition {

	var conditions []responses.Condition

	// If risk exposure is high (> 30%), suggest liquidity filter
	if riskExposure > 30 {
		conditions = append(conditions, responses.Condition{
			ID:          "min_liquidity",
			Label:       "Token liquidity above $100k",
			Description: "Only get notified about tokens that have enough money in the market to buy and sell easily. Low liquidity tokens are risky because prices can swing wildly.",
			Type:        "number",
			Field:       "liquidity",
			Operator:    "gte",
			Value:       100000,
		})
	}

	// If entry timing is high (> 70, meaning very early entries), suggest pair age filter
	// Early entries are risky — suggest waiting for token to stabilize
	if entryTiming > 70 {
		conditions = append(conditions, responses.Condition{
			ID:          "min_pair_age",
			Label:       "Token pair older than 6 hours",
			Description: "Only get notified about tokens that have been trading for at least 6 hours. Brand new tokens are more likely to be scams or crash quickly.",
			Type:        "number",
			Field:       "pair_age_hours",
			Operator:    "gte",
			Value:       6,
		})
	}

	// If token quality is low (< 60), suggest market cap filter
	if tokenQuality < 60 {
		conditions = append(conditions, responses.Condition{
			ID:          "min_mcap",
			Label:       "Market cap above $500k",
			Description: "Only get notified about tokens worth at least $500k total. Bigger tokens are generally safer and less likely to disappear overnight.",
			Type:        "number",
			Field:       "market_cap",
			Operator:    "gte",
			Value:       500000,
		})
	}

	// If win rate is below 60, suggest only following profitable token types
	// by requiring minimum 24h volume
	if len(tokenData) > 0 {
		conditions = append(conditions, responses.Condition{
			ID:          "min_volume",
			Label:       "24h trading volume above $50k",
			Description: "Only get notified about tokens that people are actively trading. Low volume means fewer buyers and sellers, making it harder to exit your position.",
			Type:        "number",
			Field:       "volume_h24",
			Operator:    "gte",
			Value:       50000,
		})
	}

	return conditions

}
