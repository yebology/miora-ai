package services

import (
	"miora-ai/app/entities"
)

// calculateMetrics computes wallet scoring based on transaction data.
// TODO: Implement real scoring logic with actual on-chain data analysis.
// For now, uses placeholder calculations based on transaction count.
func (s *WalletService) calculateMetrics(walletID uint, txs []entities.Transaction) *entities.WalletMetric {
	total := len(txs)

	// Placeholder scoring — will be replaced with real analysis
	profitConsistency := clamp(float64(total) * 1.5)
	winRate := clamp(float64(total) * 1.2)
	riskExposure := clamp(float64(total) * 0.5)
	entryTiming := clamp(float64(total) * 1.0)
	tokenQuality := clamp(float64(total) * 0.8)
	tradeDiscipline := clamp(float64(total) * 0.9)

	// Score formula from README
	finalScore := (0.25 * profitConsistency) +
		(0.20 * winRate) +
		(0.15 * entryTiming) +
		(0.15 * tokenQuality) +
		(0.15 * tradeDiscipline) -
		(0.10 * riskExposure)

	recommendation := scoreToRecommendation(finalScore)

	return &entities.WalletMetric{
		WalletID:          walletID,
		TotalTransactions: total,
		ProfitConsistency: profitConsistency,
		WinRate:           winRate,
		RiskExposure:      riskExposure,
		EntryTiming:       entryTiming,
		TokenQuality:      tokenQuality,
		TradeDiscipline:   tradeDiscipline,
		FinalScore:        finalScore,
		Recommendation:    recommendation,
	}
}

func scoreToRecommendation(score float64) string {
	switch {
	case score >= 80:
		return "full_follow"
	case score >= 60:
		return "partial_follow"
	case score >= 40:
		return "conditional_follow"
	default:
		return "avoid"
	}
}

func clamp(val float64) float64 {
	if val > 100 {
		return 100
	}
	if val < 0 {
		return 0
	}
	return val
}
