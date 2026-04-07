// Package migrations handles database schema management.
//
// This package provides three core operations:
//   - RunMigrations: auto-migrate all entity tables
//   - Reset: drop all tables and re-apply migrations (destructive)
//   - Seed: populate the database with initial development data
//
// Reset and Seed are invoked via standalone commands in cmd/reset and cmd/seed.
package migrations

import (
	"miora-ai/app/entities"

	"gorm.io/gorm"
)

// RunMigrations applies GORM auto-migration for all registered entities.
// Creates tables if they don't exist, and adds missing columns/indexes.
// Does not drop existing columns or data.
func RunMigrations(db *gorm.DB) error {

	return db.AutoMigrate(
		&entities.User{},
		&entities.Wallet{},
		&entities.Transaction{},
		&entities.WalletMetric{},
		&entities.Watchlist{},
		&entities.Notification{},
	)

}
