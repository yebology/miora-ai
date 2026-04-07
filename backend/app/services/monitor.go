// monitor.go runs a background job that polls watched wallets for new trades.
// When a new trade is detected, it checks user conditions and sends notifications
// via WebSocket (real-time) and email (async).
//
// Helper methods (poll, checkWallet, notifyFollowers, etc.) are in monitor_helper.go.
package services

import (
	"log"
	"time"

	"miora-ai/app/interfaces"
	"miora-ai/app/ws"
)

// MonitorService polls watched wallets and notifies users of new trades.
type MonitorService struct {
	watchlistRepo interfaces.IWatchlistRepository
	notifRepo     interfaces.INotificationRepository
	userRepo      interfaces.IUserRepository
	evmClient     interfaces.BlockchainClient
	svmClient     interfaces.BlockchainClient
	dexScreener   interfaces.IDexScreener
	ai            *AIService
	emailClient   interfaces.IEmailClient
	hub           *ws.Hub
	interval      time.Duration
	lastTxCount   map[string]int
}

// NewMonitorService creates a new MonitorService.
func NewMonitorService(
	watchlistRepo interfaces.IWatchlistRepository,
	notifRepo interfaces.INotificationRepository,
	userRepo interfaces.IUserRepository,
	evmClient interfaces.BlockchainClient,
	svmClient interfaces.BlockchainClient,
	dexScreener interfaces.IDexScreener,
	ai *AIService,
	emailClient interfaces.IEmailClient,
	hub *ws.Hub,
) *MonitorService {

	return &MonitorService{
		watchlistRepo: watchlistRepo,
		notifRepo:     notifRepo,
		userRepo:      userRepo,
		evmClient:     evmClient,
		svmClient:     svmClient,
		dexScreener:   dexScreener,
		ai:            ai,
		emailClient:   emailClient,
		hub:           hub,
		interval:      30 * time.Second,
		lastTxCount:   make(map[string]int),
	}

}

// Start begins the background polling loop. Call this as a goroutine.
func (m *MonitorService) Start() {

	log.Println("Monitor: started")

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for range ticker.C {
		m.poll()
	}

}
