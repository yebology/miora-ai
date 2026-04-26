// Package entities contains the User database model.
//
// Users are identified by their wallet address (MetaMask).
// Email is optional — only used for trade notifications.
package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	WalletAddress string         `gorm:"uniqueIndex;not null" json:"wallet_address"`
	Email         string         `json:"email,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
