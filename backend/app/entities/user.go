package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	FirebaseUID string         `gorm:"uniqueIndex;not null" json:"firebase_uid"`
	Email       string         `gorm:"uniqueIndex;not null" json:"email"`
	Name        string         `json:"name"`
	Avatar      string         `json:"avatar"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
