package entities

import (
	"time"

	"gorm.io/datatypes"
)

// Watchlist represents a user following a wallet with optional notification conditions.
type Watchlist struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"index;not null" json:"user_id"`
	WalletAddress  string         `gorm:"not null" json:"wallet_address"`
	Chain          string         `gorm:"not null" json:"chain"`
	Recommendation string         `json:"recommendation"`
	Conditions     datatypes.JSON `json:"conditions"` // JSON array of selected condition IDs
	EmailNotify    bool           `gorm:"default:true" json:"email_notify"`
	CreatedAt      time.Time      `json:"created_at"`
}
