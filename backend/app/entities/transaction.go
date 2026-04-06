package entities

import "time"

type Transaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	WalletID        uint      `gorm:"index;not null" json:"wallet_id"`
	Hash            string    `gorm:"uniqueIndex;not null" json:"hash"`
	Chain           string    `gorm:"not null" json:"chain"`
	From            string    `json:"from"`
	To              string    `json:"to"`
	Value           string    `json:"value"`
	TokenSymbol     string    `json:"token_symbol"`
	ContractAddress string    `json:"contract_address"`
	Direction       string    `json:"direction"` // "in" = buy/receive, "out" = sell/send
	BlockNumber     uint64    `json:"block_number"`
	Timestamp       time.Time `json:"timestamp"`
	CreatedAt       time.Time `json:"created_at"`
}
