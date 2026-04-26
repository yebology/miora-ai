package entities

import (
	"time"

	"gorm.io/datatypes"
)

type WalletMetric struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	WalletID          uint           `gorm:"uniqueIndex;not null" json:"wallet_id"`
	TotalTransactions int            `json:"total_transactions"`
	ProfitConsistency float64        `json:"profit_consistency"`
	WinRate           float64        `json:"win_rate"`
	RiskExposure      float64        `json:"risk_exposure"`
	EntryTiming       float64        `json:"entry_timing"`
	TokenQuality      float64        `json:"token_quality"`
	TradeDiscipline   float64        `json:"trade_discipline"`
	FinalScore        float64        `json:"final_score"`
	Recommendation    string         `json:"recommendation"`
	AiInsight         string         `json:"ai_insight,omitempty"`          // Stored AI insight text
	AiInsightTone     string         `json:"ai_insight_tone,omitempty"`     // Tone used: simple, eli5, custom
	AiInsightPrompt   string         `json:"ai_insight_prompt,omitempty"`   // Custom prompt (only if tone=custom)
	Conditions        datatypes.JSON `json:"conditions,omitempty"`          // Frozen conditions from analyze
	TradedTokens      datatypes.JSON `json:"traded_tokens,omitempty"`       // Frozen traded tokens from analyze
	AttestationUID    string         `json:"attestation_uid,omitempty"`     // EAS attestation UID on Base Sepolia
	AttestationTxHash string         `json:"attestation_tx_hash,omitempty"` // Transaction hash of the attestation
	UpdatedAt         time.Time      `json:"updated_at"`
}
