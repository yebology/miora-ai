// Package entities contains the AgentConfig database model.
//
// AgentConfig stores the user's AI trading agent configuration:
// budget, risk tolerance, conditions, and current status.
package entities

import (
	"time"

	"gorm.io/datatypes"
)

// AgentConfig represents a user's AI trading agent configuration.
type AgentConfig struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	UserID             uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	Budget             float64        `json:"budget"`                       // Total budget in USD
	MaxPerTrade        float64        `json:"max_per_trade"`                // Max USD per single trade
	RiskTolerance      string         `json:"risk_tolerance"`               // "low", "medium", "high"
	MinScore           int            `json:"min_score"`                    // Minimum wallet score to follow (0-100)
	Conditions         datatypes.JSON `json:"conditions"`                   // JSON array of condition IDs (same as watchlist)
	Status             string         `gorm:"default:paused" json:"status"` // "active", "paused", "stopped"
	AgentWalletAddress string         `json:"agent_wallet_address"`         // Agentic wallet address on Base Sepolia
	TotalSpent         float64        `json:"total_spent"`                  // Total USD spent by agent
	TotalTrades        int            `json:"total_trades"`                 // Number of trades executed
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}
