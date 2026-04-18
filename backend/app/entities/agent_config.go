// Package entities contains the AgentConfig database model.
//
// Two bot types:
//   - "wallet" — targets one specific wallet, copies its trades
//   - "consensus" — scans all Miora wallets, trades when multiple buy same token
package entities

import (
	"time"

	"gorm.io/datatypes"
)

// AgentConfig represents a trading bot.
type AgentConfig struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	UserID              uint           `gorm:"index;not null" json:"user_id"`
	BotType             string         `gorm:"not null;default:wallet" json:"bot_type"` // "wallet" or "consensus"
	TargetWalletAddress string         `json:"target_wallet_address,omitempty"`         // Only for wallet bots
	TargetWalletChain   string         `json:"target_wallet_chain,omitempty"`
	TargetWalletScore   int            `json:"target_wallet_score,omitempty"`
	Recommendation      string         `json:"recommendation,omitempty"`
	Budget              float64        `json:"budget"`
	MaxPerTrade         float64        `json:"max_per_trade"`
	Conditions          datatypes.JSON `json:"conditions"`
	Status              string         `gorm:"default:paused" json:"status"`
	AgentWalletAddress  string         `json:"agent_wallet_address"`
	TotalSpent          float64        `json:"total_spent"`
	TotalTrades         int            `json:"total_trades"`
	ConsensusThreshold  int            `gorm:"default:3" json:"consensus_threshold,omitempty"`   // Only for consensus bots
	ConsensusWindowMin  int            `gorm:"default:60" json:"consensus_window_min,omitempty"` // Only for consensus bots
	MinScore            int            `gorm:"default:70" json:"min_score,omitempty"`            // Only for consensus bots
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}
