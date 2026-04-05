package main

import (
	"log"

	"miora-ai/app/clients"
	"miora-ai/app/handlers"
	"miora-ai/app/repositories"
	"miora-ai/app/routes"
	"miora-ai/app/services"
	"miora-ai/config"
	"miora-ai/migrations"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	if err := migrations.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Clients
	evmClient := clients.NewAlchemyEVM(cfg.AlchemyAPIKey)
	svmClient := clients.NewAlchemySolana(cfg.AlchemyAPIKey)

	// Repository
	walletRepo := repositories.NewWalletRepository(db)

	// Service
	walletService := services.NewWalletService(walletRepo, evmClient, svmClient)

	// Handler
	walletHandler := handlers.NewWalletHandler(walletService)

	// Fiber
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	routes.SetupRoutes(app, walletHandler)

	log.Fatal(app.Listen(":" + cfg.AppPort))
}
