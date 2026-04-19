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
	"strings"
	"time"

	"miora-ai/app/clients"
	"miora-ai/app/dto"
	"miora-ai/app/entities"
	"miora-ai/app/interfaces"
)

// AgentLoopService runs the background agent trading loop.
type AgentLoopService struct {
	agentRepo   interfaces.IAgentRepository
	walletRepo  interfaces.IWalletRepository
	userRepo    interfaces.IUserRepository
	evmClient   interfaces.BlockchainClient
	dexScreener interfaces.IDexScreener
	ai          *AIService
	agentKit    *clients.AgentKitClient
	interval    time.Duration
	lastTxCount map[string]int
}

// NewAgentLoopService creates a new AgentLoopService.
func NewAgentLoopService(
	agentRepo interfaces.IAgentRepository,
	walletRepo interfaces.IWalletRepository,
	userRepo interfaces.IUserRepository,
	evmClient interfaces.BlockchainClient,
	dexScreener interfaces.IDexScreener,
	ai *AIService,
	agentKit *clients.AgentKitClient,
) *AgentLoopService {
	return &AgentLoopService{
		agentRepo:   agentRepo,
		walletRepo:  walletRepo,
		userRepo:    userRepo,
		evmClient:   evmClient,
		dexScreener: dexScreener,
		ai:          ai,
		agentKit:    agentKit,
		interval:    30 * time.Second,
		lastTxCount: make(map[string]int),
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

// processConfig evaluates trades for a single bot.
// Routes to the correct handler based on bot type.
func (s *AgentLoopService) processConfig(config *entities.AgentConfig) {
	remaining := config.Budget - config.TotalSpent
	if remaining < config.MaxPerTrade {
		return
	}

	switch config.BotType {
	case "consensus":
		s.processConsensus(config)
	default: // "wallet"
		w := walletWithScore{
			Address: config.TargetWalletAddress,
			Chain:   config.TargetWalletChain,
			Score:   config.TargetWalletScore,
		}
		s.checkWalletForAgent(config, w)
	}
}

// processConsensus scans all wallets in Miora DB and detects when multiple
// high-score wallets buy the same token within a time window.
// This is a premium feature — trades with higher confidence.
func (s *AgentLoopService) processConsensus(config *entities.AgentConfig) {
	// Get all wallets with score >= minScore from entire database
	wallets, err := s.walletRepo.FindAllWithMetrics(config.MinScore)
	if err != nil || len(wallets) == 0 {
		return
	}

	// Track which tokens are being bought by which wallets
	// Key: token contract address, Value: list of wallets that bought it
	tokenBuyers := make(map[string][]walletWithScore)

	for _, wallet := range wallets {
		metric, err := s.walletRepo.GetMetric(wallet.ID)
		if err != nil || metric == nil {
			continue
		}

		w := walletWithScore{
			Address: wallet.Address,
			Chain:   wallet.Chain,
			Score:   int(metric.FinalScore),
		}

		// Check for new buys
		transfers, err := s.evmClient.GetTransfers(wallet.Address, 10, wallet.Chain)
		if err != nil {
			continue
		}

		key := fmt.Sprintf("consensus:%d:%s", config.ID, wallet.Address)
		prevCount := s.lastTxCount[key]

		if prevCount == 0 {
			s.lastTxCount[key] = len(transfers)
			continue
		}

		if len(transfers) <= prevCount {
			continue
		}

		newTxs := transfers[:len(transfers)-prevCount]
		s.lastTxCount[key] = len(transfers)

		for _, tx := range newTxs {
			if tx.ContractAddress == "" || tx.Direction != "in" {
				continue
			}
			tokenBuyers[tx.ContractAddress] = append(tokenBuyers[tx.ContractAddress], w)
		}
	}

	// Check consensus: which tokens have >= threshold buyers?
	for tokenAddr, buyers := range tokenBuyers {
		if len(buyers) < config.ConsensusThreshold {
			continue
		}

		// Consensus reached! Multiple wallets bought the same token
		avgScore := 0
		walletList := ""
		for _, b := range buyers {
			avgScore += b.Score
			walletList += fmt.Sprintf("%s(score:%d) ", b.Address[:8], b.Score)
		}
		avgScore /= len(buyers)

		log.Printf("[AgentLoop] CONSENSUS: %d wallets bought token %s (avg score %d) — %s",
			len(buyers), tokenAddr, avgScore, walletList)

		// Use the first buyer as the "source" for the trade record
		source := buyers[0]

		// Build a synthetic transfer for the trade
		syntheticTx := interfaces.TransferData{
			ContractAddress: tokenAddr,
			TokenSymbol:     fmt.Sprintf("TOKEN_%s", tokenAddr[:8]),
			Direction:       "in",
		}

		// Execute with consensus reason
		s.evaluateAndExecuteConsensus(config, source, syntheticTx, len(buyers), avgScore)
	}
}

// evaluateAndExecuteConsensus executes a consensus-detected trade.
func (s *AgentLoopService) evaluateAndExecuteConsensus(config *entities.AgentConfig, source walletWithScore, tx interfaces.TransferData, numBuyers, avgScore int) {
	dexChain := chainToDexScreenerID(source.Chain)
	var tokenInfo *dto.TokenPairData
	pairs, err := s.dexScreener.GetTokenPairs(dexChain, tx.ContractAddress)
	if err == nil && len(pairs) > 0 {
		tokenInfo = &pairs[0]
		tx.TokenSymbol = pairs[0].BaseSymbol
	}

	if !s.meetsAgentConditions(config, tokenInfo) {
		s.recordTrade(config, source, tx, "buy", "skipped",
			fmt.Sprintf("consensus (%d wallets) but conditions not met", numBuyers), "")
		return
	}

	remaining := config.Budget - config.TotalSpent
	if config.MaxPerTrade > remaining {
		s.recordTrade(config, source, tx, "buy", "skipped", "insufficient budget", "")
		return
	}

	if !s.agentKit.IsHealthy() {
		s.recordTrade(config, source, tx, "buy", "failed", "agent sidecar unavailable", "")
		return
	}

	amountETH := fmt.Sprintf("%.6f", config.MaxPerTrade/2000)
	result, err := s.agentKit.ExecuteSwap(tx.ContractAddress, tx.TokenSymbol, amountETH, "buy")
	if err != nil {
		s.recordTrade(config, source, tx, "buy", "failed", fmt.Sprintf("swap failed: %v", err), "")
		return
	}

	log.Printf("[AgentLoop] CONSENSUS TRADE: User %d bought %s — %d wallets (avg score %d)",
		config.UserID, tx.TokenSymbol, numBuyers, avgScore)

	config.TotalSpent += config.MaxPerTrade
	config.TotalTrades++
	s.agentRepo.UpdateConfig(config)

	s.recordTrade(config, source, tx, "buy", "executed",
		fmt.Sprintf("CONSENSUS: %d wallets (avg score %d) bought %s. Agent wallet: %s",
			numBuyers, avgScore, tx.TokenSymbol, result.AgentWallet), "")
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
	direction := "buy"
	if tx.Direction == "out" {
		direction = "sell"
	}

	dexChain := chainToDexScreenerID(wallet.Chain)
	var tokenInfo *dto.TokenPairData
	pairs, err := s.dexScreener.GetTokenPairs(dexChain, tx.ContractAddress)
	if err == nil && len(pairs) > 0 {
		tokenInfo = &pairs[0]
	}

	// For buys: check conditions and budget
	if direction == "buy" {
		if !s.meetsAgentConditions(config, tokenInfo) {
			s.recordTrade(config, wallet, tx, direction, "skipped", "conditions not met", "")
			return
		}

		remaining := config.Budget - config.TotalSpent
		tradeAmount := config.MaxPerTrade
		if tradeAmount > remaining {
			s.recordTrade(config, wallet, tx, direction, "skipped", "insufficient budget", "")
			return
		}
	}

	// AI risk assessment
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

	if !s.agentKit.IsHealthy() {
		s.recordTrade(config, wallet, tx, direction, "failed", "agent sidecar unavailable", riskAssessment)
		return
	}

	amountETH := fmt.Sprintf("%.6f", config.MaxPerTrade/2000)
	result, err := s.agentKit.ExecuteSwap(tx.ContractAddress, tx.TokenSymbol, amountETH, direction)
	if err != nil {
		s.recordTrade(config, wallet, tx, direction, "failed", fmt.Sprintf("swap failed: %v", err), riskAssessment)
		return
	}

	log.Printf("[AgentLoop] User %d: %s %s (%.2f USD) triggered by wallet %s (score %d)",
		config.UserID, direction, tx.TokenSymbol, config.MaxPerTrade, wallet.Address, wallet.Score)

	if direction == "buy" {
		config.TotalSpent += config.MaxPerTrade
	} else {
		// Auto-transfer sell proceeds to user's connected wallet
		user, err := s.userRepo.FindByID(config.UserID)
		if err == nil && user.WalletAddress != "" {
			amountToTransfer := fmt.Sprintf("%.6f", config.MaxPerTrade/2000)
			if _, transferErr := s.agentKit.ExecuteTransfer(user.WalletAddress, amountToTransfer); transferErr != nil {
				log.Printf("[AgentLoop] Auto-transfer to %s failed: %v", user.WalletAddress, transferErr)
			} else {
				log.Printf("[AgentLoop] Auto-transferred %s ETH to user %s", amountToTransfer, user.WalletAddress)
			}
		}
	}
	config.TotalTrades++
	s.agentRepo.UpdateConfig(config)

	s.recordTrade(config, wallet, tx, direction, "executed",
		fmt.Sprintf("%s %s because wallet %s (score %d) did the same. Agent wallet: %s",
			strings.ToUpper(direction[:1])+direction[1:], tx.TokenSymbol, wallet.Address, wallet.Score, result.AgentWallet),
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
func (s *AgentLoopService) recordTrade(config *entities.AgentConfig, wallet walletWithScore, tx interfaces.TransferData, direction, status, reason, riskAssessment string) {
	s.agentRepo.CreateTrade(&entities.AgentTrade{
		AgentConfigID:  config.ID,
		SourceWallet:   wallet.Address,
		SourceScore:    wallet.Score,
		TokenAddress:   tx.ContractAddress,
		TokenSymbol:    tx.TokenSymbol,
		Direction:      direction,
		AmountUSD:      config.MaxPerTrade,
		Status:         status,
		Reason:         reason,
		RiskAssessment: riskAssessment,
	})
}
