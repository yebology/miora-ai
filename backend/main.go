// Miora AI — Backend Entry Point
//
// This is the main entry point for the Miora AI backend server.
// It initializes the application in the following order:
//
//  1. Load environment configuration from .env
//  2. Connect to PostgreSQL database via GORM
//  3. Run auto-migrations for all entities
//  4. Set up Fiber HTTP server with middleware
//  5. Wire all dependencies and register routes via router.SetUp()
//
// The server listens on the port defined by APP_PORT in .env.
package main

import (
	"log"

	"miora-ai/app/ws"
	"miora-ai/config"
	"miora-ai/migrations"
	"miora-ai/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// Load environment variables and validate required keys
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Establish database connection using DSN from config
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run auto-migrations for Wallet, Transaction, and WalletMetric tables
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Fiber with request logger and CORS middleware
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:  cfg.AllowedOrigins,
		AllowMethods:  "*",
		AllowHeaders:  "Origin,Content-Type,Accept,Authorization,X-Wallet-Address",
		ExposeHeaders: "Content-Length",
	}))

	// Initialize WebSocket hub
	hub := ws.NewHub()

	// Wire dependencies and register routes (also starts monitor)
	router.SetUp(app, db, cfg, hub)

	log.Fatal(app.Listen(":" + cfg.AppPort))

}
