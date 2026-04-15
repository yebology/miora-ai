// Package interfaces defines contracts for the AI trading agent.
package interfaces

import (
	"miora-ai/app/entities"
	"miora-ai/pkg"
)

// IAgentService defines the business logic contract for the AI trading agent.
type IAgentService interface {
	// GetOrCreateConfig returns the agent config for a user, creating a default if none exists.
	GetOrCreateConfig(userID uint) (*entities.AgentConfig, *pkg.AppError)

	// UpdateConfig updates the agent configuration for a user.
	UpdateConfig(userID uint, budget, maxPerTrade float64, riskTolerance string, minScore int, conditions []string) (*entities.AgentConfig, *pkg.AppError)

	// Start activates the agent for a user.
	Start(userID uint) (*entities.AgentConfig, *pkg.AppError)

	// Pause pauses the agent for a user.
	Pause(userID uint) (*entities.AgentConfig, *pkg.AppError)

	// GetStatus returns the current agent status and summary.
	GetStatus(userID uint) (*entities.AgentConfig, *pkg.AppError)

	// GetTrades returns the agent's trade history for a user.
	GetTrades(userID uint, limit int) ([]entities.AgentTrade, *pkg.AppError)
}

// IAgentRepository defines the data access contract for agent operations.
type IAgentRepository interface {
	FindConfigByUserID(userID uint) (*entities.AgentConfig, error)
	FindActiveConfigs() ([]entities.AgentConfig, error)
	CreateConfig(config *entities.AgentConfig) error
	UpdateConfig(config *entities.AgentConfig) error
	FindTradesByConfigID(configID uint, limit int) ([]entities.AgentTrade, error)
	CreateTrade(trade *entities.AgentTrade) error
}
