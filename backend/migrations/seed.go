package migrations

import (
	"log"
	"miora-ai/app/entities"

	"gorm.io/gorm"
)

// Seed populates the database with initial development data.
// Uses FirstOrCreate to skip records that already exist (based on address unique index).
// Safe to run multiple times without duplicating data.
func Seed(db *gorm.DB) error {

	wallets := []entities.Wallet{
		{Address: "0x0000000000000000000000000000000000000001", Chain: "evm"},
		{Address: "0x0000000000000000000000000000000000000002", Chain: "evm"},
	}

	for _, w := range wallets {
		result := db.FirstOrCreate(&w, entities.Wallet{Address: w.Address})
		if result.Error != nil {
			return result.Error
		}
	}

	log.Println("Seed: done")
	return nil

}
