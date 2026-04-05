package migrations

import (
	"miora-ai/app/entities"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.Wallet{},
		&entities.Transaction{},
		&entities.WalletMetric{},
	)
}
