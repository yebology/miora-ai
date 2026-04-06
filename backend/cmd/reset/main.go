// Command reset drops all database tables and re-applies migrations.
// This is a destructive operation intended for development use only.
//
// Usage:
//
//	go run cmd/reset/main.go
//	make db-reset
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

	if err := migrations.Reset(db); err != nil {
		log.Fatalf("Reset failed: %v", err)
	}

}
