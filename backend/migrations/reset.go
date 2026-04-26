package migrations

import (
	"log"
	"miora-ai/app/entities"

	"gorm.io/gorm"
)

// Reset drops all entity tables in reverse dependency order,
// then re-applies migrations from scratch.
//
// Drop order matters — WalletMetric and Transaction reference Wallet,
// so they must be dropped first to avoid foreign key violations.
//
// WARNING: This destroys all existing data. Use only in development.
func Reset(db *gorm.DB) error {

	tables := []interface{}{
		&entities.WalletMetric{},
		&entities.Transaction{},
		&entities.Wallet{},
	}

	for _, t := range tables {
		if err := db.Migrator().DropTable(t); err != nil {
			return err
		}
	}

	log.Println("Reset: tables dropped")

	if err := RunMigrations(db); err != nil {
		return err
	}

	log.Println("Reset: migrations re-applied")
	return nil

}
