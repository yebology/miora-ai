// Scoring logic for Miora AI wallet intelligence.
//
// # Overview
//
// This file calculates a wallet's "intelligence score" based on its trading behavior.
// The score tells users whether a wallet is worth following or not.
//
// # Data Sources
//
// Two external APIs provide the data needed for scoring:
//
// 1. DexScreener (free, no API key):
//   - Liquidity: how much money is available to trade a token.
//     Low liquidity (< $10k) means the token is risky — hard to buy/sell without
//     moving the price. Used to calculate Risk Exposure.
//   - Market Cap: total value of all tokens in circulation.
//     Higher market cap = more established token. Used to calculate Token Quality.
//   - Pair Created At: when the trading pair was first created on a DEX.
//     If a wallet trades tokens shortly after creation, it means early entry.
//     Used to calculate Entry Timing.
//   - Price Change 24h: how much the token price changed in the last 24 hours.
//     Used as a fallback for Win Rate when Moralis data is not available.
//
// 2. Moralis (requires API key, free tier available):
//   - Historical Token Price: the USD price of a token at a specific block number.
//     This lets us know the price when the wallet bought the token.
//   - Current Token Price: the USD price of a token right now.
//   - PnL (Profit and Loss): calculated as ((exit - buy) / buy) × 100.
//     Example: bought at $1.00, sold at $1.50 → PnL = +50%.
//     Used to calculate Win Rate and Profit Consistency with real data.
//
// # Metrics (each scored 0–100)
//
//   - Win Rate: percentage of trades where the wallet made a profit.
//     If Moralis data is available: uses actual PnL per trade (realized + unrealized).
//     If not: uses DexScreener 24h price change as approximation.
//     Example: 4 trades, 3 profitable → winRate = 75.
//
//   - Profit Consistency: how stable the profits are across trades.
//     Calculated using standard deviation (stdDev) of PnL values.
//     stdDev measures how "spread out" the numbers are from the average.
//     Low stdDev = consistent profits = high score.
//     High stdDev = volatile/unpredictable results = low score.
//     Formula: profitConsistency = 100 - stdDev (clamped 0–100).
//     Example: PnL = [+20%, -5%, +15%, +10%] → stdDev ≈ 9.35 → score = 90.65.
//     Example: PnL = [+80%, -50%, +90%, -40%] → stdDev ≈ 61 → score = 38.97.
//
//   - Risk Exposure: percentage of traded tokens with liquidity below $10,000.
//     Low liquidity tokens are risky — easy to manipulate, hard to exit.
//     Example: 10 tokens traded, 3 have liquidity < $10k → riskExposure = 30.
//
//   - Entry Timing: how early the wallet enters new tokens.
//     Based on the average age of trading pairs when the wallet traded them.
//     Younger pairs = earlier entry = higher score (sniper behavior).
//     Score 100 if average pair age < 24 hours, scales linearly to 0 at 720 hours (30 days).
//     Example: avg pair age = 48 hours → score = 100 - (48/720 × 100) = 93.33.
//     Example: avg pair age = 360 hours → score = 100 - (360/720 × 100) = 50.
//
//   - Token Quality: average market cap of tokens the wallet trades.
//     Uses logarithmic scale (log10) because market caps range from $100 to $100B.
//     Without log scale, one $10B token would dominate the average.
//     Score = (log10(avgMcap) / 7) × 100, where 10^7 = $10M = score 100.
//     Example: avg mcap = $10M → log10(10M) = 7 → score = 100.
//     Example: avg mcap = $100k → log10(100k) = 5 → score = 71.43.
//     Example: avg mcap = $1k → log10(1k) = 3 → score = 42.86.
//
//   - Trade Discipline: ratio of unique tokens traded vs total transactions.
//     A focused wallet trades few tokens many times (low ratio = high score).
//     A scattered wallet trades many tokens once each (high ratio = low score).
//     Formula: tradeDiscipline = (1 - (uniqueTokens / totalTxs)) × 100.
//     Example: 5 tokens across 100 txs → ratio = 0.05 → score = 95.
//     Example: 80 tokens across 100 txs → ratio = 0.80 → score = 20.
//
// # Final Score Formula
//
// Weights sum to 1.0 for a true 0–100 range:
//
//	Final Score = (0.30 × Profit Consistency)
//	            + (0.30 × Win Rate)
//	            + (0.15 × Entry Timing)
//	            + (0.15 × Token Quality)
//	            + (0.10 × Trade Discipline)
//
// Risk Exposure is calculated and returned in the response for informational purposes,
// but is NOT included in the final score formula.
//
// # Recommendation
//
//	80–100 → full_follow    (consistent, profitable, low risk)
//	40–79  → partial_follow (mixed results, follow with caution)
//	< 40   → avoid          (high risk, poor performance)
package services

import (
	"math"
	"time"

	"miora-ai/app/dto"
	"miora-ai/app/entities"
	"miora-ai/utils"
)

// calculateMetrics computes wallet scoring from transaction data,
// DexScreener token data, and Moralis PnL data.
//
// Parameters:
//   - walletID: database ID of the wallet being analyzed
//   - txs: all transactions (buy + sell) fetched from Alchemy
//   - tokenData: DexScreener data per unique token (liquidity, mcap, pair age, price change)
//   - trades: PnL results per trade from Moralis (buy price vs exit price)
func (s *WalletService) calculateMetrics(
	walletID uint,
	txs []entities.Transaction,
	tokenData map[string]dto.TokenPairData,
	trades []tradeResult,
) *entities.WalletMetric {

	total := len(txs)
	if total == 0 {
		return &entities.WalletMetric{
			WalletID:       walletID,
			Recommendation: "avoid",
		}
	}

	// Safety: prevent divide-by-zero in calculations below.
	// If no token data from DexScreener, set to 1 so divisions return 0 instead of crashing.
	tokenCount := float64(len(tokenData))
	if tokenCount == 0 {
		tokenCount = 1
	}

	// --- Win Rate & Profit Consistency ---
	// Two paths depending on whether Moralis PnL data is available.
	var winRate, profitConsistency float64

	if len(trades) > 0 {
		// PRIMARY: Real PnL from Moralis (more accurate)
		// wins = number of trades with positive PnL
		// winRate = (wins / totalTrades) × 100
		//
		// profitConsistency uses standard deviation:
		//   1. Sum all PnL values (sumPnl) and their squares (sumPnlSq)
		//   2. mean = sumPnl / count
		//   3. variance = (sumPnlSq / count) - (mean²)
		//      → measures how spread out the PnL values are
		//   4. stdDev = √variance
		//      → converts back to same unit (percentage)
		//   5. profitConsistency = 100 - stdDev
		//      → lower spread = higher consistency score
		wins := 0.0
		var sumPnl, sumPnlSq float64
		tradeCount := float64(len(trades))

		for _, t := range trades {
			if t.PnlPercent > 0 {
				wins++
			}
			sumPnl += t.PnlPercent
			sumPnlSq += t.PnlPercent * t.PnlPercent
		}

		winRate = utils.Clamp((wins / tradeCount) * 100)

		mean := sumPnl / tradeCount
		variance := (sumPnlSq / tradeCount) - (mean * mean)
		stdDev := math.Sqrt(math.Abs(variance))
		profitConsistency = utils.Clamp(100 - stdDev)

	} else {
		// FALLBACK: DexScreener 24h price change (less accurate)
		// Used when Moralis data is unavailable (e.g. Solana, or API errors).
		// Same math as above, but using price change instead of actual PnL.
		wins := 0.0
		var sumChange, sumChangeSq float64

		for _, t := range tokenData {
			if t.PriceChangeH24 > 0 {
				wins++
			}
			sumChange += t.PriceChangeH24
			sumChangeSq += t.PriceChangeH24 * t.PriceChangeH24
		}

		winRate = utils.Clamp((wins / tokenCount) * 100)

		mean := sumChange / tokenCount
		variance := (sumChangeSq / tokenCount) - (mean * mean)
		stdDev := math.Sqrt(math.Abs(variance))
		profitConsistency = utils.Clamp(100 - stdDev)
	}

	// --- Risk Exposure ---
	// Calculated and stored for informational purposes only — NOT included in final score.
	// Count how many tokens have liquidity below the configured threshold.
	// Threshold is configurable via SCORING_LIQUIDITY_THRESHOLD in .env.
	lowLiq := 0.0
	for _, t := range tokenData {
		if t.Liquidity < s.scoring.LiquidityThreshold {
			lowLiq++
		}
	}
	riskExposure := utils.Clamp((lowLiq / tokenCount) * 100)

	// --- Entry Timing ---
	// Calculate average pair age in hours for all tokens the wallet traded.
	// PairCreatedAt is in unix milliseconds from DexScreener.
	// Convert to hours: (now_ms - pairCreatedAt_ms) / 3,600,000.
	// Score scales linearly: 0 hours = 100, maxAge hours = 0.
	// Max age is configurable via SCORING_ENTRY_TIMING_MAX_AGE in .env.
	// Default 50 if no pair age data is available.
	now := float64(time.Now().UnixMilli())
	var totalAge float64
	ageCount := 0.0
	for _, t := range tokenData {
		if t.PairCreatedAt > 0 {
			ageHours := (now - float64(t.PairCreatedAt)) / 3600000
			totalAge += ageHours
			ageCount++
		}
	}
	entryTiming := 50.0
	if ageCount > 0 {
		avgAge := totalAge / ageCount
		entryTiming = utils.Clamp(100 - (avgAge / s.scoring.EntryTimingMaxAge * 100))
	}

	// --- Token Quality ---
	// Average market cap of all tokens traded, scored on a logarithmic scale.
	// Log scale is used because market caps range wildly ($100 to $100B).
	// Log base is configurable via SCORING_TOKEN_QUALITY_LOG_BASE in .env.
	// Default 7 means: log10($10M) = 7 = score 100.
	var totalMcap float64
	for _, t := range tokenData {
		totalMcap += t.MarketCap
	}
	avgMcap := totalMcap / tokenCount
	tokenQuality := utils.Clamp((math.Log10(math.Max(avgMcap, 1)) / s.scoring.TokenQualityLogBase) * 100)

	// --- Trade Discipline ---
	// Measures how focused the wallet is.
	// ratio = uniqueTokens / totalTransactions.
	// Low ratio = trades few tokens repeatedly = disciplined/focused.
	// High ratio = trades many different tokens = scattered/unfocused.
	// Score = (1 - ratio) × 100.
	ratio := tokenCount / float64(total)
	tradeDiscipline := utils.Clamp((1 - ratio) * 100)

	// --- Final Score ---
	// Weighted sum of all metrics. Weights sum to 1.0 for a true 0–100 range.
	// Risk Exposure is NOT included in the formula — informational only.
	finalScore := (0.30 * profitConsistency) +
		(0.30 * winRate) +
		(0.15 * entryTiming) +
		(0.15 * tokenQuality) +
		(0.10 * tradeDiscipline)

	recommendation := scoreToRecommendation(finalScore)

	return &entities.WalletMetric{
		WalletID:          walletID,
		TotalTransactions: total,
		ProfitConsistency: utils.Round2(profitConsistency),
		WinRate:           utils.Round2(winRate),
		RiskExposure:      utils.Round2(riskExposure),
		EntryTiming:       utils.Round2(entryTiming),
		TokenQuality:      utils.Round2(tokenQuality),
		TradeDiscipline:   utils.Round2(tradeDiscipline),
		FinalScore:        utils.Round2(finalScore),
		Recommendation:    recommendation,
	}

}

// scoreToRecommendation converts a final score (0–100) to a recommendation label.
//
//	80–100 → "full_follow"
//	40–79  → "partial_follow"
//	< 40   → "avoid"
func scoreToRecommendation(score float64) string {

	switch {
	case score >= 80:
		return "full_follow"
	case score >= 40:
		return "partial_follow"
	default:
		return "avoid"
	}

}
