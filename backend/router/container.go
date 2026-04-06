// Package router handles dependency injection and route registration.
//
// container.go is the central DI container. It wires all application
// dependencies in the correct order:
//
//	Clients → Repositories → Services → Handlers
//
// Each layer only depends on the layer below it.
// Services depend on repositories via interfaces.
// Handlers depend on services via interfaces.
//
// To add a new domain (e.g. "dashboard"):
//  1. Create repository: dashboardRepo := repositories.NewDashboardRepository(db)
//  2. Create service:    dashboardService := services.NewDashboardService(dashboardRepo)
//  3. Create handler:    dashboardHandler := handlers.NewDashboardHandler(dashboardService)
//  4. Add to Container struct and return
package router

import (
	"miora-ai/app/clients"
	"miora-ai/app/handlers"
	"miora-ai/app/repositories"
	"miora-ai/app/services"
	"miora-ai/config"

	"gorm.io/gorm"
)

// Container holds all handler instances, ready to be used by route registration.
type Container struct {
	WalletHandler *handlers.WalletHandler
}

// NewContainer creates all dependencies and returns a fully wired Container.
// Initialization order: clients → repositories → services → handlers.
func NewContainer(db *gorm.DB, alchemyAPIKey, moralisAPIKey, birdeyeAPIKey string, scoring config.ScoringConfig) *Container {

	// Clients
	evmClient := clients.NewAlchemyEVM(alchemyAPIKey)
	svmClient := clients.NewAlchemySolana(alchemyAPIKey)
	dexScreener := clients.NewDexScreener()
	moralis := clients.NewMoralis(moralisAPIKey)
	birdeye := clients.NewBirdeye(birdeyeAPIKey)

	// Repositories
	walletRepo := repositories.NewWalletRepository(db)

	// Services
	walletService := services.NewWalletService(walletRepo, evmClient, svmClient, dexScreener, moralis, birdeye, scoring)

	// Handlers
	walletHandler := handlers.NewWalletHandler(walletService)

	return &Container{
		WalletHandler: walletHandler,
	}

}
