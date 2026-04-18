// Package interfaces defines contracts for the AI trading agent.
package interfaces

import (
	"miora-ai/app/entities"
	"miora-ai/pkg"
)

// IAgentService defines the business logic contract for the AI trading agent.
type IAgentService interface {
	// CreateBot creates a new bot.
	CreateBot(userID uint, botType, targetWallet, chain string, score int, recommendation string, budget, maxPerTrade float64, conditions []string, consensusThreshold, consensusWindowMin, minScore int) (*entities.AgentConfig, *pkg.AppError)

	// UpdateBot updates a bot's configuration.
	UpdateBot(botID, userID uint, budget, maxPerTrade float64, conditions []string, consensusThreshold, consensusWindowMin, minScore int) (*entities.AgentConfig, *pkg.AppError)

	// DeleteBot removes a bot.
	DeleteBot(botID, userID uint) *pkg.AppError

	// StartBot activates a bot.
	StartBot(botID, userID uint) (*entities.AgentConfig, *pkg.AppError)

	// PauseBot pauses a bot.
	PauseBot(botID, userID uint) (*entities.AgentConfig, *pkg.AppError)

	// GetBot returns a single bot by ID.
	GetBot(botID, userID uint) (*entities.AgentConfig, *pkg.AppError)

	// ListBots returns all bots for a user.
	ListBots(userID uint) ([]entities.AgentConfig, *pkg.AppError)

	// GetTrades returns trade history for a bot.
	GetTrades(botID, userID uint, limit int) ([]entities.AgentTrade, *pkg.AppError)
}

// IAgentRepository defines the data access contract for agent operations.
type IAgentRepository interface {
	FindByID(id uint) (*entities.AgentConfig, error)
	FindByUserID(userID uint) ([]entities.AgentConfig, error)
	FindActiveConfigs() ([]entities.AgentConfig, error)
	CreateConfig(config *entities.AgentConfig) error
	UpdateConfig(config *entities.AgentConfig) error
	DeleteConfig(id uint) error
	FindTradesByConfigID(configID uint, limit int) ([]entities.AgentTrade, error)
	CreateTrade(trade *entities.AgentTrade) error
}
