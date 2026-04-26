// Package entities contains the AgentTrade database model.
//
// AgentTrade records each trade executed by the AI trading agent,
// including the source wallet that triggered it and the execution details.
package entities

import "time"

// AgentTrade represents a single trade executed by the AI agent.
type AgentTrade struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	AgentConfigID  uint      `gorm:"index;not null" json:"agent_config_id"`
	SourceWallet   string    `json:"source_wallet"`   // Wallet address that triggered the trade
	SourceScore    int       `json:"source_score"`    // Score of the source wallet at time of trade
	TokenAddress   string    `json:"token_address"`   // Token contract address
	TokenSymbol    string    `json:"token_symbol"`    // Token symbol (e.g. "PEPE")
	Direction      string    `json:"direction"`       // "buy" or "sell"
	AmountUSD      float64   `json:"amount_usd"`      // Amount in USD
	TxHash         string    `json:"tx_hash"`         // Transaction hash on Base Sepolia
	Status         string    `json:"status"`          // "executed", "failed", "skipped"
	Reason         string    `json:"reason"`          // Why the trade was executed/skipped
	RiskAssessment string    `json:"risk_assessment"` // AI risk assessment for this trade
	CreatedAt      time.Time `json:"created_at"`
}
