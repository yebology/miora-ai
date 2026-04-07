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
	"miora-ai/app/ws"
	"miora-ai/config"

	"gorm.io/gorm"
)

// Container holds all handler instances, ready to be used by route registration.
type Container struct {
	WalletHandler    *handlers.WalletHandler
	SwapHandler      *handlers.SwapHandler
	AuthHandler      *handlers.AuthHandler
	WatchlistHandler *handlers.WatchlistHandler
	Monitor          *services.MonitorService
}

// NewContainer creates all dependencies and returns a fully wired Container.
// Initialization order: clients → repositories → services → handlers.
func NewContainer(db *gorm.DB, alchemyAPIKey, moralisAPIKey, birdeyeAPIKey, geminiAPIKey, oneInchAPIKey string, scoring config.ScoringConfig, hub *ws.Hub) *Container {

	// Clients
	evmClient := clients.NewAlchemyEVM(alchemyAPIKey)
	svmClient := clients.NewAlchemySolana(alchemyAPIKey)
	dexScreener := clients.NewDexScreener()
	moralis := clients.NewMoralis(moralisAPIKey)
	birdeye := clients.NewBirdeye(birdeyeAPIKey)
	gemini := clients.NewGemini(geminiAPIKey)
	jupiter := clients.NewJupiter()
	oneInch := clients.NewOneInch(oneInchAPIKey)

	// Repositories
	walletRepo := repositories.NewWalletRepository(db)
	userRepo := repositories.NewUserRepository(db)
	watchlistRepo := repositories.NewWatchlistRepository(db)
	notifRepo := repositories.NewNotificationRepository(db)

	// Services
	aiService := services.NewAIService(gemini)
	walletService := services.NewWalletService(walletRepo, evmClient, svmClient, dexScreener, moralis, birdeye, aiService, scoring)
	swapService := services.NewSwapService(jupiter, oneInch)
	userService := services.NewUserService(userRepo)
	watchlistService := services.NewWatchlistService(watchlistRepo)

	// Monitor
	monitorService := services.NewMonitorService(watchlistRepo, notifRepo, evmClient, svmClient, dexScreener, aiService, hub)

	// Handlers
	walletHandler := handlers.NewWalletHandler(walletService)
	swapHandler := handlers.NewSwapHandler(swapService)
	authHandler := handlers.NewAuthHandler(userService)
	watchlistHandler := handlers.NewWatchlistHandler(watchlistService, userService)

	return &Container{
		WalletHandler:    walletHandler,
		SwapHandler:      swapHandler,
		AuthHandler:      authHandler,
		WatchlistHandler: watchlistHandler,
		Monitor:          monitorService,
	}

}
