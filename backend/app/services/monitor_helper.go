// monitor_helper.go contains helper methods for MonitorService.
// These handle polling logic, wallet checking, notification dispatch,
// condition evaluation, and safe token data access.
package services

import (
	"encoding/json"
	"log"
	"miora-ai/app/dto"
	"miora-ai/app/entities"
	"miora-ai/app/interfaces"
	"miora-ai/app/ws"
	"miora-ai/constants"
	"time"
)

// poll checks all unique watched wallets for new transactions.
func (m *MonitorService) poll() {

	wallets := m.getUniqueWatchedWallets()

	for _, w := range wallets {
		m.checkWallet(w.WalletAddress, w.Chain)
	}

}

// getUniqueWatchedWallets returns deduplicated wallet addresses from all watchlists.
func (m *MonitorService) getUniqueWatchedWallets() []entities.Watchlist {

	// Get all watchlist items (across all users)
	// In production, this should be optimized with a distinct query
	allItems, err := m.watchlistRepo.FindByWallet("")
	if err != nil {
		return nil
	}

	seen := make(map[string]bool)
	var unique []entities.Watchlist
	for _, item := range allItems {
		key := item.WalletAddress + ":" + item.Chain
		if !seen[key] {
			seen[key] = true
			unique = append(unique, item)
		}
	}

	return unique

}

// checkWallet fetches latest transactions and detects new ones.
func (m *MonitorService) checkWallet(address, chain string) {

	if !constants.IsEVM(chain) {
		return
	}

	client := m.evmClient

	transfers, err := client.GetTransfers(address, 100, chain)
	if err != nil {
		return
	}

	key := address + ":" + chain
	prevCount := m.lastTxCount[key]

	if prevCount == 0 {
		// First poll — just record count, don't notify
		m.lastTxCount[key] = len(transfers)
		return
	}

	if len(transfers) <= prevCount {
		return // No new transactions
	}

	// New transactions detected
	newTxs := transfers[:len(transfers)-prevCount]
	m.lastTxCount[key] = len(transfers)

	log.Printf("Monitor: %d new txs for %s on %s", len(newTxs), address, chain)

	// For each new tx with a token contract, notify followers
	for _, tx := range newTxs {
		if tx.ContractAddress == "" {
			continue // Skip native transfers without token contract
		}

		m.notifyFollowers(address, chain, tx)
	}

}

// notifyFollowers sends notifications to all users following this wallet.
func (m *MonitorService) notifyFollowers(walletAddress, chain string, tx interfaces.TransferData) {

	followers, err := m.watchlistRepo.FindByWallet(walletAddress)
	if err != nil || len(followers) == 0 {
		return
	}

	// Fetch token data from DexScreener
	dexChain := chainToDexScreenerID(chain)
	var tokenInfo *dto.TokenPairData
	pairs, err := m.dexScreener.GetTokenPairs(dexChain, tx.ContractAddress)
	if err == nil && len(pairs) > 0 {
		tokenInfo = &pairs[0]
	}

	// Generate AI risk assessment for this trade (once, reuse for all followers)
	aiAssessment := ""
	if m.ai != nil && tokenInfo != nil {
		pairAgeHours := float64(0)
		if tokenInfo.PairCreatedAt > 0 {
			pairAgeHours = float64(time.Now().UnixMilli()-tokenInfo.PairCreatedAt) / 3600000
		}
		assessment, err := m.ai.GenerateTradeAssessment(
			walletAddress, chain, tx.TokenSymbol, tx.Direction,
			tokenInfo.Liquidity, tokenInfo.MarketCap, tokenInfo.PriceChangeH24, pairAgeHours,
		)
		if err == nil {
			aiAssessment = assessment
		} else {
			log.Printf("Monitor: AI assessment failed for %s: %v", tx.TokenSymbol, err)
		}
	}

	for _, follower := range followers {
		// Check if token meets user's conditions
		if !m.meetsConditions(follower, tokenInfo) {
			continue
		}

		// Build notification
		notification := ws.Message{
			Type: "wallet_trade",
			Payload: map[string]any{
				"wallet_address":   walletAddress,
				"chain":            chain,
				"token_address":    tx.ContractAddress,
				"token_symbol":     tx.TokenSymbol,
				"direction":        tx.Direction,
				"value":            tx.Value,
				"traded_at":        time.Now(),
				"liquidity":        getTokenLiquidity(tokenInfo),
				"market_cap":       getTokenMcap(tokenInfo),
				"price_change_24h": getTokenPriceChange(tokenInfo),
				"ai_assessment":    aiAssessment,
			},
		}

		// Send via WebSocket
		m.hub.SendToUser(follower.UserID, notification)

		// Save to database for history
		m.notifRepo.Create(&entities.Notification{
			UserID:        follower.UserID,
			WalletAddress: walletAddress,
			Chain:         chain,
			TokenAddress:  tx.ContractAddress,
			TokenSymbol:   tx.TokenSymbol,
			Direction:     tx.Direction,
			Value:         tx.Value,
			Liquidity:     getTokenLiquidity(tokenInfo),
			MarketCap:     getTokenMcap(tokenInfo),
			AiAssessment:  aiAssessment,
		})

		log.Printf("Monitor: notified user %d about %s trading %s", follower.UserID, walletAddress, tx.TokenSymbol)
	}

}

// meetsConditions checks if a token meets the user's selected conditions.
func (m *MonitorService) meetsConditions(follower entities.Watchlist, tokenInfo *dto.TokenPairData) bool {

	if tokenInfo == nil {
		return true // No data to check — let it through
	}

	// Parse selected condition IDs
	var selectedConditions []string
	json.Unmarshal(follower.Conditions, &selectedConditions)

	if len(selectedConditions) == 0 {
		return true // No conditions set — notify always
	}

	for _, cond := range selectedConditions {
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

// Helper functions for safe token data access

func getTokenLiquidity(t *dto.TokenPairData) float64 {

	if t != nil {
		return t.Liquidity
	}
	return 0

}

func getTokenMcap(t *dto.TokenPairData) float64 {

	if t != nil {
		return t.MarketCap
	}
	return 0

}

func getTokenPriceChange(t *dto.TokenPairData) float64 {

	if t != nil {
		return t.PriceChangeH24
	}
	return 0

}
