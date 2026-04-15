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
	"log"
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
	WalletHandler     *handlers.WalletHandler
	SwapHandler       *handlers.SwapHandler
	AuthHandler       *handlers.AuthHandler
	WatchlistHandler  *handlers.WatchlistHandler
	ReputationHandler *handlers.ReputationHandler
	AgentHandler      *handlers.AgentHandler
	Monitor           *services.MonitorService
	AgentLoop         *services.AgentLoopService
}

// NewContainer creates all dependencies and returns a fully wired Container.
// Initialization order: clients → repositories → services → handlers.
func NewContainer(db *gorm.DB, alchemyAPIKey, moralisAPIKey, geminiAPIKey, oneInchAPIKey, resendAPIKey, resendFrom string, scoring config.ScoringConfig, easCfg config.EASConfig, hub *ws.Hub) *Container {

	// Clients
	evmClient := clients.NewAlchemyEVM(alchemyAPIKey)
	dexScreener := clients.NewDexScreener()
	moralis := clients.NewMoralis(moralisAPIKey)
	gemini := clients.NewGemini(geminiAPIKey)
	oneInch := clients.NewOneInch(oneInchAPIKey)

	// EAS Client (Base Sepolia)
	easClient, err := clients.NewEASClient(easCfg.RPCURL, easCfg.EASContractAddress, easCfg.SchemaUID, easCfg.AttesterPrivateKey)
	if err != nil {
		log.Printf("[WARN] EAS client initialization failed: %v — attestations will be disabled", err)
		easClient = &clients.EASClient{}
	}

	// Repositories
	walletRepo := repositories.NewWalletRepository(db)
	userRepo := repositories.NewUserRepository(db)
	watchlistRepo := repositories.NewWatchlistRepository(db)
	notifRepo := repositories.NewNotificationRepository(db)
	agentRepo := repositories.NewAgentRepository(db)

	// Services
	aiService := services.NewAIService(gemini)
	walletService := services.NewWalletService(walletRepo, evmClient, dexScreener, moralis, aiService, scoring, easClient)
	swapService := services.NewSwapService(oneInch)
	userService := services.NewUserService(userRepo)
	watchlistService := services.NewWatchlistService(watchlistRepo)
	agentService := services.NewAgentService(agentRepo)

	// Monitor
	resendClient := clients.NewResend(resendAPIKey, resendFrom)
	monitorService := services.NewMonitorService(watchlistRepo, notifRepo, userRepo, evmClient, dexScreener, aiService, resendClient, hub)

	// Agent Loop (background trading)
	agentKitClient := clients.NewAgentKitClient("")
	agentLoop := services.NewAgentLoopService(agentRepo, walletRepo, evmClient, dexScreener, aiService, agentKitClient)

	// Handlers
	walletHandler := handlers.NewWalletHandler(walletService)
	swapHandler := handlers.NewSwapHandler(swapService)
	authHandler := handlers.NewAuthHandler(userService)
	watchlistHandler := handlers.NewWatchlistHandler(watchlistService, userService)
	reputationHandler := handlers.NewReputationHandler(walletRepo, easClient)
	agentHandler := handlers.NewAgentHandler(agentService, userService)

	return &Container{
		WalletHandler:     walletHandler,
		SwapHandler:       swapHandler,
		AuthHandler:       authHandler,
		WatchlistHandler:  watchlistHandler,
		ReputationHandler: reputationHandler,
		AgentHandler:      agentHandler,
		Monitor:           monitorService,
		AgentLoop:         agentLoop,
	}

}
