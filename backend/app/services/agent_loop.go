// Package services contains the AI trading agent background loop.
//
// AgentLoopService runs as a goroutine, polling for active agent configs
// and monitoring wallets from the user's watchlist for new trades.
// When a trade is detected that meets the agent's conditions, it executes
// a swap via the AgentKit Python sidecar.
//
// Flow:
//  1. Poll every 30 seconds
//  2. Get all active agent configs
//  3. For each config, get user's watchlist (wallets they already follow)
//  4. Check each watched wallet for new trades
//  5. Evaluate: conditions met? budget available? AI risk ok?
//  6. If all pass → call Python sidecar to execute swap
//  7. Record trade in database
package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"miora-ai/app/clients"
	"miora-ai/app/dto"
	"miora-ai/app/entities"
	"miora-ai/app/interfaces"
)

// AgentLoopService runs the background agent trading loop.
type AgentLoopService struct {
	agentRepo     interfaces.IAgentRepository
	watchlistRepo interfaces.IWatchlistRepository
	walletRepo    interfaces.IWalletRepository
	evmClient     interfaces.BlockchainClient
	dexScreener   interfaces.IDexScreener
	ai            *AIService
	agentKit      *clients.AgentKitClient
	interval      time.Duration
	lastTxCount   map[string]int
}

// NewAgentLoopService creates a new AgentLoopService.
func NewAgentLoopService(
	agentRepo interfaces.IAgentRepository,
	watchlistRepo interfaces.IWatchlistRepository,
	walletRepo interfaces.IWalletRepository,
	evmClient interfaces.BlockchainClient,
	dexScreener interfaces.IDexScreener,
	ai *AIService,
	agentKit *clients.AgentKitClient,
) *AgentLoopService {
	return &AgentLoopService{
		agentRepo:     agentRepo,
		watchlistRepo: watchlistRepo,
		walletRepo:    walletRepo,
		evmClient:     evmClient,
		dexScreener:   dexScreener,
		ai:            ai,
		agentKit:      agentKit,
		interval:      30 * time.Second,
		lastTxCount:   make(map[string]int),
	}
}

// Start begins the background agent loop. Call as a goroutine.
func (s *AgentLoopService) Start() {
	log.Println("[AgentLoop] Started")

	if !s.agentKit.IsHealthy() {
		log.Println("[AgentLoop] WARNING: AgentKit sidecar not available — agent trading disabled")
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for range ticker.C {
		s.poll()
	}
}

// poll checks all active agent configs and evaluates trades.
func (s *AgentLoopService) poll() {
	configs, err := s.agentRepo.FindActiveConfigs()
	if err != nil || len(configs) == 0 {
		return
	}

	for i := range configs {
		s.processConfig(&configs[i])
	}
}

// walletWithScore holds a wallet address with its score for agent evaluation.
type walletWithScore struct {
	Address string
	Chain   string
	Score   int
}

// processConfig evaluates trades for a single agent config.
// Only monitors wallets from the user's watchlist — not all wallets in the database.
func (s *AgentLoopService) processConfig(config *entities.AgentConfig) {
	// Check budget
	remaining := config.Budget - config.TotalSpent
	if remaining < config.MaxPerTrade {
		return
	}

	// Get wallets from user's watchlist
	watchlist, err := s.watchlistRepo.FindByUser(config.UserID)
	if err != nil || len(watchlist) == 0 {
		return
	}

	// For each watched wallet, get its score and check for new trades
	for _, item := range watchlist {
		wallet, err := s.walletRepo.FindByAddress(item.WalletAddress)
		if err != nil || wallet == nil {
			continue
		}

		metric, err := s.walletRepo.GetMetric(wallet.ID)
		if err != nil || metric == nil {
			continue
		}

		// Skip wallets below minScore
		if int(metric.FinalScore) < config.MinScore {
			continue
		}

		w := walletWithScore{
			Address: item.WalletAddress,
			Chain:   item.Chain,
			Score:   int(metric.FinalScore),
		}

		s.checkWalletForAgent(config, w)
	}
}

// checkWalletForAgent checks a wallet for new trades and evaluates them for the agent.
func (s *AgentLoopService) checkWalletForAgent(config *entities.AgentConfig, wallet walletWithScore) {
	transfers, err := s.evmClient.GetTransfers(wallet.Address, 10, wallet.Chain)
	if err != nil {
		return
	}

	key := fmt.Sprintf("agent:%d:%s:%s", config.ID, wallet.Address, wallet.Chain)
	prevCount := s.lastTxCount[key]

	if prevCount == 0 {
		s.lastTxCount[key] = len(transfers)
		return
	}

	if len(transfers) <= prevCount {
		return
	}

	newTxs := transfers[:len(transfers)-prevCount]
	s.lastTxCount[key] = len(transfers)

	for _, tx := range newTxs {
		if tx.ContractAddress == "" {
			continue
		}
		s.evaluateAndExecute(config, wallet, tx)
	}
}

// evaluateAndExecute evaluates a trade and executes it if all conditions pass.
func (s *AgentLoopService) evaluateAndExecute(config *entities.AgentConfig, wallet walletWithScore, tx interfaces.TransferData) {
	if tx.Direction != "in" {
		return
	}

	dexChain := chainToDexScreenerID(wallet.Chain)
	var tokenInfo *dto.TokenPairData
	pairs, err := s.dexScreener.GetTokenPairs(dexChain, tx.ContractAddress)
	if err == nil && len(pairs) > 0 {
		tokenInfo = &pairs[0]
	}

	if !s.meetsAgentConditions(config, tokenInfo) {
		s.recordTrade(config, wallet, tx, "skipped", "conditions not met", "")
		return
	}

	remaining := config.Budget - config.TotalSpent
	tradeAmount := config.MaxPerTrade
	if tradeAmount > remaining {
		s.recordTrade(config, wallet, tx, "skipped", "insufficient budget", "")
		return
	}

	riskAssessment := ""
	if s.ai != nil && tokenInfo != nil {
		pairAgeHours := float64(0)
		if tokenInfo.PairCreatedAt > 0 {
			pairAgeHours = float64(time.Now().UnixMilli()-tokenInfo.PairCreatedAt) / 3600000
		}
		assessment, err := s.ai.GenerateTradeAssessment(
			wallet.Address, wallet.Chain, tx.TokenSymbol, tx.Direction,
			tokenInfo.Liquidity, tokenInfo.MarketCap, tokenInfo.PriceChangeH24, pairAgeHours,
		)
		if err == nil {
			riskAssessment = assessment
		}
	}

	if config.RiskTolerance == "low" && tokenInfo != nil && tokenInfo.Liquidity < 100000 {
		s.recordTrade(config, wallet, tx, "skipped", "risk too high for low tolerance", riskAssessment)
		return
	}

	if !s.agentKit.IsHealthy() {
		s.recordTrade(config, wallet, tx, "failed", "agent sidecar unavailable", riskAssessment)
		return
	}

	amountETH := fmt.Sprintf("%.6f", tradeAmount/2000)
	result, err := s.agentKit.ExecuteSwap(tx.ContractAddress, tx.TokenSymbol, amountETH)
	if err != nil {
		s.recordTrade(config, wallet, tx, "failed", fmt.Sprintf("swap failed: %v", err), riskAssessment)
		return
	}

	log.Printf("[AgentLoop] User %d: executed swap for %s (%.2f USD) triggered by wallet %s (score %d)",
		config.UserID, tx.TokenSymbol, tradeAmount, wallet.Address, wallet.Score)

	config.TotalSpent += tradeAmount
	config.TotalTrades++
	s.agentRepo.UpdateConfig(config)

	s.recordTrade(config, wallet, tx, "executed",
		fmt.Sprintf("Bought %s because wallet %s (score %d) bought it. Agent wallet: %s",
			tx.TokenSymbol, wallet.Address, wallet.Score, result.AgentWallet),
		riskAssessment)
}

// meetsAgentConditions checks if a token meets the agent's conditions.
func (s *AgentLoopService) meetsAgentConditions(config *entities.AgentConfig, tokenInfo *dto.TokenPairData) bool {
	if tokenInfo == nil {
		return false
	}

	var conditions []string
	json.Unmarshal(config.Conditions, &conditions)

	if len(conditions) == 0 {
		return true
	}

	for _, cond := range conditions {
		switch cond {
		case "min_liquidity":
			if tokenInfo.Liquidity < 100000 {
				return false
			}
		case "min_pair_age":
			if tokenInfo.PairCreatedAt > 0 {
				ageHours := float64(time.Now().UnixMilli()-tokenInfo.PairCreatedAt) / 3600000
				if ageHours < 6 {
					return false
				}
			}
		case "min_mcap":
			if tokenInfo.MarketCap < 500000 {
				return false
			}
		case "min_volume":
			if tokenInfo.VolumeH24 < 50000 {
				return false
			}
		}
	}

	return true
}

// recordTrade saves an agent trade to the database.
func (s *AgentLoopService) recordTrade(config *entities.AgentConfig, wallet walletWithScore, tx interfaces.TransferData, status, reason, riskAssessment string) {
	s.agentRepo.CreateTrade(&entities.AgentTrade{
		AgentConfigID:  config.ID,
		SourceWallet:   wallet.Address,
		SourceScore:    wallet.Score,
		TokenAddress:   tx.ContractAddress,
		TokenSymbol:    tx.TokenSymbol,
		Direction:      "buy",
		AmountUSD:      config.MaxPerTrade,
		Status:         status,
		Reason:         reason,
		RiskAssessment: riskAssessment,
	})
}
