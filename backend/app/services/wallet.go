// Package services contains business logic for each domain.
//
// Each service implements an interface from app/interfaces/.
// Services depend on repository interfaces and blockchain client interfaces.
// Methods return *pkg.AppError for structured error handling in the handler layer.
//
// Files in this package:
//   - wallet.go:        main service logic (AnalyzeWallet, GetWallet, calculateTrades)
//   - wallet_helper.go: helper methods (findOrCreateWallet, getPrice, fetchTokenData, etc.)
//   - scoring.go:       scoring calculation logic
//   - ai.go:            AI insight generation
//   - swap.go:          swap quote service
package services

import (
	"log"
	"time"

	"miora-ai/app/dto/responses"
	"miora-ai/app/entities"
	"miora-ai/app/interfaces"
	"miora-ai/config"
	"miora-ai/constants"
	"miora-ai/pkg"
)

// WalletService implements interfaces.IWalletService.
type WalletService struct {
	repo        interfaces.IWalletRepository
	evmClient   interfaces.BlockchainClient
	dexScreener interfaces.IDexScreener
	moralis     interfaces.IMoralis
	ai          *AIService
	scoring     config.ScoringConfig
	eas         interfaces.IEASClient
}

// NewWalletService creates a new WalletService with the given dependencies.
func NewWalletService(
	repo interfaces.IWalletRepository,
	evmClient interfaces.BlockchainClient,
	dexScreener interfaces.IDexScreener,
	moralis interfaces.IMoralis,
	ai *AIService,
	scoring config.ScoringConfig,
	eas interfaces.IEASClient,
) *WalletService {

	return &WalletService{
		repo:        repo,
		evmClient:   evmClient,
		dexScreener: dexScreener,
		moralis:     moralis,
		ai:          ai,
		scoring:     scoring,
		eas:         eas,
	}

}

// tradeResult holds PnL data for a single trade (buy → sell/hold).
type tradeResult struct {
	TokenAddress string
	BuyPrice     float64
	ExitPrice    float64   // sell price (realized) or current price (unrealized)
	PnlPercent   float64   // ((exit - buy) / buy) * 100
	BuyTime      time.Time // when the wallet bought
	ExitTime     time.Time // when the wallet sold (zero if unrealized)
}

// AnalyzeWallet orchestrates the full analysis flow and returns scoring.
func (s *WalletService) AnalyzeWallet(address, chain string, limit int) (*responses.WalletAnalysis, *pkg.AppError) {

	if !constants.IsValidTransactionLimit(chain, limit) {
		limit = constants.GetTransactionLimits(chain).Default
	}

	wallet, appErr := s.findOrCreateWallet(address, chain)
	if appErr != nil {
		return nil, appErr
	}

	txEntities, appErr := s.fetchAndSaveTransactions(wallet, chain, limit)
	if appErr != nil {
		return nil, appErr
	}

	tokenData := s.fetchTokenData(chain, txEntities)
	trades := s.calculateTrades(chain, txEntities)

	metric := s.calculateMetrics(wallet.ID, txEntities, tokenData, trades)
	if err := s.repo.SaveMetric(metric); err != nil {
		return nil, pkg.ErrInternal()
	}

	result := &responses.WalletAnalysis{
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
		TradedTokens:      buildTradedTokens(chain, trades, txEntities),
	}

	// Generate conditions for conditional_follow recommendations
	if metric.Recommendation == "conditional_follow" {
		result.Conditions = buildConditions(
			tokenData,
			metric.RiskExposure, metric.EntryTiming, metric.TokenQuality,
			s.scoring,
		)
	}

	// Generate AI insight (non-blocking — if it fails, return without insight)
	if insight, err := s.ai.GenerateInsight(result, "simple"); err == nil {
		result.AiInsight = insight
	} else {
		log.Printf("AI insight failed: %v", err)
	}

	// Publish EAS attestation on Base Sepolia (non-blocking — if it fails, continue without attestation)
	if s.eas != nil {
		go func() {
			uid, txHash, err := s.eas.Attest(
				address,
				uint8(metric.FinalScore),
				metric.Recommendation,
				uint32(metric.TotalTransactions),
				chain,
			)
			if err != nil {
				log.Printf("[EAS] Attestation failed for %s: %v", address, err)
				return
			}
			log.Printf("[EAS] Attestation published for %s — UID: %s, TxHash: %s", address, uid, txHash)

			// Update metric with attestation data
			metric.AttestationUID = uid
			metric.AttestationTxHash = txHash
			if saveErr := s.repo.SaveMetric(metric); saveErr != nil {
				log.Printf("[EAS] Failed to save attestation UID for %s: %v", address, saveErr)
			}
		}()
	}

	return result, nil

}

// GetWallet retrieves a previously analyzed wallet by address.
func (s *WalletService) GetWallet(address string) (*responses.WalletAnalysis, *pkg.AppError) {

	wallet, err := s.repo.FindByAddress(address)
	if err != nil {
		return nil, pkg.ErrNotFound(constants.DataNotFound)
	}

	metric, err := s.repo.GetMetric(wallet.ID)
	if err != nil {
		return nil, pkg.ErrNotFound(constants.DataNotFound)
	}

	return &responses.WalletAnalysis{
		Address:           wallet.Address,
		Chain:             wallet.Chain,
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

// RegenerateInsight regenerates the AI insight for a previously analyzed wallet with a different tone.
func (s *WalletService) RegenerateInsight(address, chain, tone, customPrompt string) (string, *pkg.AppError) {

	result, appErr := s.GetWallet(address)
	if appErr != nil {
		return "", appErr
	}
	result.Chain = chain

	var insight string
	var err error

	if tone == "custom" && customPrompt != "" {
		insight, err = s.ai.GenerateCustomInsight(result, customPrompt)
	} else {
		insight, err = s.ai.GenerateInsight(result, tone)
	}

	if err != nil {
		return "", pkg.ErrInternal()
	}

	return insight, nil

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
			var exitTime time.Time

			// Try to match with next sell (FIFO)
			if sellIdx < len(group.sells) {
				sell := group.sells[sellIdx]
				price := s.getPrice(chain, addr, sell.BlockNumber, sell.Timestamp)
				if price > 0 {
					exitPrice = price
					exitTime = sell.Timestamp
					sellIdx++
				}
			}

			// No sell → unrealized, use current price
			if exitPrice == 0 {
				exitPrice = s.getPrice(chain, addr, 0, time.Now())
				if exitPrice == 0 {
					continue
				}
				// exitTime stays zero — indicates unrealized
			}

			pnl := ((exitPrice - buyPrice) / buyPrice) * 100
			results = append(results, tradeResult{
				TokenAddress: addr,
				BuyPrice:     buyPrice,
				ExitPrice:    exitPrice,
				PnlPercent:   pnl,
				BuyTime:      buy.Timestamp,
				ExitTime:     exitTime,
			})
		}
	}

	return results

}
