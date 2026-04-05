package entities

import "time"

type WalletMetric struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	WalletID          uint      `gorm:"uniqueIndex;not null" json:"wallet_id"`
	TotalTransactions int       `json:"total_transactions"`
	ProfitConsistency float64   `json:"profit_consistency"`
	WinRate           float64   `json:"win_rate"`
	RiskExposure      float64   `json:"risk_exposure"`
	EntryTiming       float64   `json:"entry_timing"`
	TokenQuality      float64   `json:"token_quality"`
	TradeDiscipline   float64   `json:"trade_discipline"`
	FinalScore        float64   `json:"final_score"`
	Recommendation    string    `json:"recommendation"`
	UpdatedAt         time.Time `json:"updated_at"`
}
