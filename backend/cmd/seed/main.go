// Command seed populates the database with initial development data.
// Safe to run multiple times — existing records are skipped.
//
// Usage:
//
//	go run cmd/seed/main.go
//	make db-seed
package main

import (
	"log"

	"miora-ai/config"
	"miora-ai/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := migrations.Seed(db); err != nil {
		log.Fatalf("Seed failed: %v", err)
	}

}
