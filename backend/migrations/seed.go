package migrations

import (
	"encoding/json"
	"log"
	"time"

	"miora-ai/app/entities"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Seed populates the database with demo data for hackathon presentation.
// Uses FirstOrCreate to skip records that already exist.
// Safe to run multiple times without duplicating data.
func Seed(db *gorm.DB) error {

	// --- 1. Demo user (wallet-based auth) ---
	// Replace this address with your actual MetaMask wallet address before demo.
	demoUser := entities.User{
		WalletAddress: "0xe5bd375e6f4650c13740b81a4e95f3b01836b5e0",
	}
	db.FirstOrCreate(&demoUser, entities.User{WalletAddress: demoUser.WalletAddress})
	log.Printf("Seed: user ID=%d wallet=%s", demoUser.ID, demoUser.WalletAddress)

	// --- 2. Analyzed wallets with metrics ---
	targetWallet := entities.Wallet{Address: "0x28c6c06298d514db089934071355e5743bf21d60", Chain: "base"}
	db.FirstOrCreate(&targetWallet, entities.Wallet{Address: targetWallet.Address})

	metric := entities.WalletMetric{
		WalletID:          targetWallet.ID,
		TotalTransactions: 47,
		ProfitConsistency: 72.5,
		WinRate:           68.0,
		RiskExposure:      15.3,
		EntryTiming:       81.2,
		TokenQuality:      65.8,
		TradeDiscipline:   77.4,
		FinalScore:        73,
		Recommendation:    "conditional_follow",
	}
	db.FirstOrCreate(&metric, entities.WalletMetric{WalletID: targetWallet.ID})

	// High-score wallet for full_follow demo
	topWallet := entities.Wallet{Address: "0x21a31ee1afc51d94c2efccaa2092ad1028285549", Chain: "base"}
	db.FirstOrCreate(&topWallet, entities.Wallet{Address: topWallet.Address})

	topMetric := entities.WalletMetric{
		WalletID:          topWallet.ID,
		TotalTransactions: 82,
		ProfitConsistency: 88.4,
		WinRate:           85.0,
		RiskExposure:      5.2,
		EntryTiming:       74.6,
		TokenQuality:      91.3,
		TradeDiscipline:   89.1,
		FinalScore:        88,
		Recommendation:    "full_follow",
	}
	db.FirstOrCreate(&topMetric, entities.WalletMetric{WalletID: topWallet.ID})

	// --- 3. Watchlist entries ---
	conditionsJSON, _ := json.Marshal([]string{"min_liquidity", "min_mcap"})
	watchlist1 := entities.Watchlist{
		UserID:         demoUser.ID,
		WalletAddress:  targetWallet.Address,
		Chain:          "base",
		Recommendation: "conditional_follow",
		Conditions:     datatypes.JSON(conditionsJSON),
		EmailNotify:    false,
	}
	db.FirstOrCreate(&watchlist1, entities.Watchlist{UserID: demoUser.ID, WalletAddress: targetWallet.Address})

	emptyConditions, _ := json.Marshal([]string{})
	watchlist2 := entities.Watchlist{
		UserID:         demoUser.ID,
		WalletAddress:  topWallet.Address,
		Chain:          "base",
		Recommendation: "full_follow",
		Conditions:     datatypes.JSON(emptyConditions),
		EmailNotify:    false,
	}
	db.FirstOrCreate(&watchlist2, entities.Watchlist{UserID: demoUser.ID, WalletAddress: topWallet.Address})

	// --- 4. Bot configs ---
	walletBot := entities.AgentConfig{
		UserID:              demoUser.ID,
		BotType:             "wallet",
		TargetWalletAddress: targetWallet.Address,
		TargetWalletChain:   "base",
		TargetWalletScore:   73,
		Recommendation:      "conditional_follow",
		Budget:              500,
		MaxPerTrade:         50,
		Conditions:          datatypes.JSON(conditionsJSON),
		Status:              "active",
		AgentWalletAddress:  "0x742d35Cc6634C0532925a3b844Bc9e7595f2bD18",
		TotalSpent:          145,
		TotalTrades:         4,
	}
	db.FirstOrCreate(&walletBot, entities.AgentConfig{UserID: demoUser.ID, BotType: "wallet", TargetWalletAddress: targetWallet.Address})

	consensusBot := entities.AgentConfig{
		UserID:             demoUser.ID,
		BotType:            "consensus",
		Budget:             300,
		MaxPerTrade:        30,
		Conditions:         datatypes.JSON(emptyConditions),
		Status:             "paused",
		AgentWalletAddress: "",
		TotalSpent:         0,
		TotalTrades:        0,
		ConsensusThreshold: 3,
		ConsensusWindowMin: 60,
		MinScore:           75,
	}
	db.FirstOrCreate(&consensusBot, entities.AgentConfig{UserID: demoUser.ID, BotType: "consensus"})

	// --- 5. Agent trades (demo trade history) ---
	now := time.Now()
	trades := []entities.AgentTrade{
		{
			AgentConfigID:  walletBot.ID,
			SourceWallet:   targetWallet.Address,
			SourceScore:    73,
			TokenAddress:   "0x6982508145454ce325ddbe47a25d4ec3d2311933",
			TokenSymbol:    "PEPE",
			Direction:      "buy",
			AmountUSD:      45,
			TxHash:         "0xabc123def456789012345678901234567890abcdef1234567890abcdef123456",
			Status:         "executed",
			Reason:         "Bought PEPE because wallet 0x28c6...1d60 (score 73) bought it. Agent wallet: 0x742d...bD18",
			RiskAssessment: "Moderate risk. Token has $180k liquidity and $2.1M market cap. Pair age 14 hours. Meets all conditions.",
			CreatedAt:      now.Add(-2 * time.Hour),
		},
		{
			AgentConfigID:  walletBot.ID,
			SourceWallet:   targetWallet.Address,
			SourceScore:    73,
			TokenAddress:   "0x514910771af9ca656af840dff83e8264ecf986ca",
			TokenSymbol:    "LINK",
			Direction:      "buy",
			AmountUSD:      50,
			TxHash:         "0xdef789abc012345678901234567890abcdef1234567890abcdef123456789012",
			Status:         "executed",
			Reason:         "Bought LINK because wallet 0x28c6...1d60 (score 73) bought it. Agent wallet: 0x742d...bD18",
			RiskAssessment: "Low risk. Token has $45M liquidity and $8.2B market cap. Well-established token.",
			CreatedAt:      now.Add(-5 * time.Hour),
		},
		{
			AgentConfigID:  walletBot.ID,
			SourceWallet:   targetWallet.Address,
			SourceScore:    73,
			TokenAddress:   "0x0000000000000000000000000000000000001337",
			TokenSymbol:    "NEWTOKEN",
			Direction:      "buy",
			AmountUSD:      50,
			TxHash:         "",
			Status:         "skipped",
			Reason:         "Token liquidity $8k — below min_liquidity condition ($100k minimum).",
			RiskAssessment: "High risk. Token has only $8k liquidity and $12k market cap. Pair created 45 minutes ago. Likely a new launch with high rug risk.",
			CreatedAt:      now.Add(-8 * time.Hour),
		},
		{
			AgentConfigID:  walletBot.ID,
			SourceWallet:   targetWallet.Address,
			SourceScore:    73,
			TokenAddress:   "0x514910771af9ca656af840dff83e8264ecf986ca",
			TokenSymbol:    "LINK",
			Direction:      "sell",
			AmountUSD:      50,
			TxHash:         "0x999888abc012345678901234567890abcdef1234567890abcdef123456789012",
			Status:         "executed",
			Reason:         "Sold LINK because wallet 0x28c6...1d60 (score 73) sold it. Proceeds auto-transferred to user wallet.",
			RiskAssessment: "",
			CreatedAt:      now.Add(-1 * time.Hour),
		},
	}

	for i := range trades {
		db.FirstOrCreate(&trades[i], entities.AgentTrade{
			AgentConfigID: trades[i].AgentConfigID,
			TokenSymbol:   trades[i].TokenSymbol,
			Direction:     trades[i].Direction,
			CreatedAt:     trades[i].CreatedAt,
		})
	}

	log.Printf("Seed: wallet bot ID=%d, consensus bot ID=%d, %d trades seeded", walletBot.ID, consensusBot.ID, len(trades))
	log.Println("Seed: done")
	return nil
}
