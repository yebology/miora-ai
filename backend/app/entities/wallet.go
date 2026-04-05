package entities

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Address   string         `gorm:"uniqueIndex;not null" json:"address"`
	Chain     string         `gorm:"not null" json:"chain"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
