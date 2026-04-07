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
	svmClient   interfaces.BlockchainClient
	dexScreener interfaces.IDexScreener
	moralis     interfaces.IMoralis
	birdeye     interfaces.IBirdeye
	ai          *AIService
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
	ai *AIService,
	scoring config.ScoringConfig,
) *WalletService {

	return &WalletService{
		repo:        repo,
		evmClient:   evmClient,
		svmClient:   svmClient,
		dexScreener: dexScreener,
		moralis:     moralis,
		birdeye:     birdeye,
		ai:          ai,
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

	// Generate AI insight (non-blocking — if it fails, return without insight)
	if insight, err := s.ai.GenerateInsight(result); err == nil {
		result.AiInsight = insight
	} else {
		log.Printf("AI insight failed: %v", err)
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
