package entities

import "time"

// Notification stores a trade alert sent to a user.
type Notification struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	WalletAddress string    `gorm:"not null" json:"wallet_address"`
	Chain         string    `json:"chain"`
	TokenAddress  string    `json:"token_address"`
	TokenSymbol   string    `json:"token_symbol"`
	Direction     string    `json:"direction"` // "in" = buy, "out" = sell
	Value         string    `json:"value"`     // Amount traded (human-readable)
	Liquidity     float64   `json:"liquidity"`
	MarketCap     float64   `json:"market_cap"`
	AiAssessment  string    `json:"ai_assessment"` // AI-generated risk assessment for this trade
	Read          bool      `gorm:"default:false" json:"read"`
	CreatedAt     time.Time `json:"created_at"`
}
