// Package services contains the AI trading agent business logic.
//
// The agent monitors top-scored wallets, evaluates their trades through
// the scoring engine and AI risk assessment, and executes swaps when
// all conditions are met.
//
// Note: For the hackathon, the agent service handles config management
// and trade recording. Actual autonomous trading (monitor → evaluate → execute)
// is a background process that will be implemented with AgentKit integration.
package services

import (
	"encoding/json"

	"miora-ai/app/entities"
	"miora-ai/app/interfaces"
	"miora-ai/constants"
	"miora-ai/pkg"
)

// AgentService implements interfaces.IAgentService.
type AgentService struct {
	repo interfaces.IAgentRepository
}

// NewAgentService creates a new AgentService.
func NewAgentService(repo interfaces.IAgentRepository) *AgentService {
	return &AgentService{repo: repo}
}

// GetOrCreateConfig returns the agent config for a user, creating a default if none exists.
func (s *AgentService) GetOrCreateConfig(userID uint) (*entities.AgentConfig, *pkg.AppError) {
	config, err := s.repo.FindConfigByUserID(userID)
	if err == nil && config != nil {
		return config, nil
	}

	// Create empty config — user must configure via PUT /agent/config before starting
	defaultConditions, _ := json.Marshal([]string{})
	config = &entities.AgentConfig{
		UserID:     userID,
		Conditions: defaultConditions,
		Status:     constants.DefaultAgentStatus,
	}

	if err := s.repo.CreateConfig(config); err != nil {
		return nil, pkg.ErrInternal()
	}

	return config, nil
}

// UpdateConfig updates the agent configuration for a user.
func (s *AgentService) UpdateConfig(userID uint, budget, maxPerTrade float64, riskTolerance string, minScore int, conditions []string) (*entities.AgentConfig, *pkg.AppError) {
	config, appErr := s.GetOrCreateConfig(userID)
	if appErr != nil {
		return nil, appErr
	}

	if budget > 0 {
		config.Budget = budget
	}
	if maxPerTrade > 0 {
		config.MaxPerTrade = maxPerTrade
	}
	if riskTolerance != "" {
		config.RiskTolerance = riskTolerance
	}
	if minScore >= 0 && minScore <= 100 {
		config.MinScore = minScore
	}
	if conditions != nil {
		condJSON, _ := json.Marshal(conditions)
		config.Conditions = condJSON
	}

	if err := s.repo.UpdateConfig(config); err != nil {
		return nil, pkg.ErrInternal()
	}

	return config, nil
}

// Start activates the agent for a user.
func (s *AgentService) Start(userID uint) (*entities.AgentConfig, *pkg.AppError) {
	config, appErr := s.GetOrCreateConfig(userID)
	if appErr != nil {
		return nil, appErr
	}

	config.Status = "active"
	if err := s.repo.UpdateConfig(config); err != nil {
		return nil, pkg.ErrInternal()
	}

	return config, nil
}

// Pause pauses the agent for a user.
func (s *AgentService) Pause(userID uint) (*entities.AgentConfig, *pkg.AppError) {
	config, appErr := s.GetOrCreateConfig(userID)
	if appErr != nil {
		return nil, appErr
	}

	config.Status = "paused"
	if err := s.repo.UpdateConfig(config); err != nil {
		return nil, pkg.ErrInternal()
	}

	return config, nil
}

// GetStatus returns the current agent status.
func (s *AgentService) GetStatus(userID uint) (*entities.AgentConfig, *pkg.AppError) {
	return s.GetOrCreateConfig(userID)
}

// GetTrades returns the agent's trade history for a user.
func (s *AgentService) GetTrades(userID uint, limit int) ([]entities.AgentTrade, *pkg.AppError) {
	config, err := s.repo.FindConfigByUserID(userID)
	if err != nil || config == nil {
		return []entities.AgentTrade{}, nil
	}

	trades, err := s.repo.FindTradesByConfigID(config.ID, limit)
	if err != nil {
		return nil, pkg.ErrInternal()
	}

	return trades, nil
}
