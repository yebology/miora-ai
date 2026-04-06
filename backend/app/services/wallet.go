// Package services contains business logic for each domain.
//
// Each service implements an interface from app/interfaces/.
// Services depend on repository interfaces and blockchain client interfaces.
// Methods return *pkg.AppError for structured error handling in the handler layer.
package services

import (
	"errors"
	"log"
	"time"

	"miora-ai/app/dto"
	"miora-ai/app/dto/responses"
	"miora-ai/app/entities"
	"miora-ai/app/interfaces"
	"miora-ai/config"
	"miora-ai/constants"
	"miora-ai/pkg"

	"gorm.io/gorm"
)

// WalletService implements interfaces.IWalletService.
type WalletService struct {
	repo        interfaces.IWalletRepository
	evmClient   interfaces.BlockchainClient
	svmClient   interfaces.BlockchainClient
	dexScreener interfaces.IDexScreener
	moralis     interfaces.IMoralis
	birdeye     interfaces.IBirdeye
	scoring     config.ScoringConfig
}

// NewWalletService creates a new WalletService with the given dependencies.
func NewWalletService(
	repo interfaces.IWalletRepository,
	evmClient interfaces.BlockchainClient,
	svmClient interfaces.BlockchainClient,
	dexScreener interfaces.IDexScreener,
	moralis interfaces.IMoralis,
	birdeye interfaces.IBirdeye,
	scoring config.ScoringConfig,
) *WalletService {

	return &WalletService{
		repo:        repo,
		evmClient:   evmClient,
		svmClient:   svmClient,
		dexScreener: dexScreener,
		moralis:     moralis,
		birdeye:     birdeye,
		scoring:     scoring,
	}

}

// tradeResult holds PnL data for a single trade (buy → sell/hold).
type tradeResult struct {
	TokenAddress string
	BuyPrice     float64
	ExitPrice    float64 // sell price (realized) or current price (unrealized)
	PnlPercent   float64 // ((exit - buy) / buy) * 100
}

// AnalyzeWallet orchestrates the full analysis flow and returns scoring.
func (s *WalletService) AnalyzeWallet(address, chain string) (*responses.WalletAnalysis, *pkg.AppError) {

	wallet, appErr := s.findOrCreateWallet(address, chain)
	if appErr != nil {
		return nil, appErr
	}

	txEntities, appErr := s.fetchAndSaveTransactions(wallet, chain)
	if appErr != nil {
		return nil, appErr
	}

	tokenData := s.fetchTokenData(chain, txEntities)
	trades := s.calculateTrades(chain, txEntities)

	metric := s.calculateMetrics(wallet.ID, txEntities, tokenData, trades)
	if err := s.repo.SaveMetric(metric); err != nil {
		return nil, pkg.ErrInternal()
	}

	return &responses.WalletAnalysis{
		Address:           address,
		Chain:             chain,
		TotalTransactions: metric.TotalTransactions,
		ProfitConsistency: metric.ProfitConsistency,
		WinRate:           metric.WinRate,
		RiskExposure:      metric.RiskExposure,
		EntryTiming:       metric.EntryTiming,
		TokenQuality:      metric.TokenQuality,
		TradeDiscipline:   metric.TradeDiscipline,
		FinalScore:        metric.FinalScore,
		Recommendation:    metric.Recommendation,
	}, nil

}

// --- Helper methods ---

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

	switch chain {
	case "evm":
		return s.evmClient, nil
	case "svm":
		return s.svmClient, nil
	default:
		return nil, pkg.ErrBadReq(constants.UnsupportedChain)
	}

}

// fetchAndSaveTransactions fetches transfers from Alchemy and persists them.
func (s *WalletService) fetchAndSaveTransactions(wallet *entities.Wallet, chain string) ([]entities.Transaction, *pkg.AppError) {

	client, appErr := s.getClient(chain)
	if appErr != nil {
		return nil, appErr
	}

	transfers, err := client.GetTransfers(wallet.Address)
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

// getPrice fetches token price — Moralis for EVM (by block), Birdeye for Solana (by timestamp).
// block=0 and timestamp=now means current price.
func (s *WalletService) getPrice(chain, tokenAddr string, blockNumber uint64, timestamp time.Time) float64 {

	if chain == "svm" {
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

// calculateTrades matches buy/sell per token (FIFO) and calculates PnL.
func (s *WalletService) calculateTrades(chain string, txs []entities.Transaction) []tradeResult {

	grouped := s.groupByToken(txs)
	var results []tradeResult

	for addr, group := range grouped {
		sellIdx := 0

		for _, buy := range group.buys {
			buyPrice := s.getPrice(chain, addr, buy.BlockNumber, buy.Timestamp)
			if buyPrice == 0 {
				log.Printf("Price: skip buy %s: no price data", addr)
				continue
			}

			exitPrice := 0.0

			// Try to match with next sell (FIFO)
			if sellIdx < len(group.sells) {
				sell := group.sells[sellIdx]
				price := s.getPrice(chain, addr, sell.BlockNumber, sell.Timestamp)
				if price > 0 {
					exitPrice = price
					sellIdx++
				}
			}

			// No sell → unrealized, use current price
			if exitPrice == 0 {
				exitPrice = s.getPrice(chain, addr, 0, time.Now())
				if exitPrice == 0 {
					continue
				}
			}

			pnl := ((exitPrice - buyPrice) / buyPrice) * 100
			results = append(results, tradeResult{
				TokenAddress: addr,
				BuyPrice:     buyPrice,
				ExitPrice:    exitPrice,
				PnlPercent:   pnl,
			})
		}
	}

	return results

}

// --- Utility methods ---

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

func chainToDexScreenerID(chain string) string {

	switch chain {
	case "evm":
		return "ethereum"
	case "svm":
		return "solana"
	default:
		return chain
	}

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
