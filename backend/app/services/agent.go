// Package services contains the AI trading agent business logic.
//
// Each bot targets one specific wallet. A user can have multiple bots.
// Conditions are inherited from the wallet's analyze result.
package services

import (
	"encoding/json"

	"miora-ai/app/entities"
	"miora-ai/app/interfaces"
	"miora-ai/constants"
	"miora-ai/pkg"
)

type AgentService struct {
	repo interfaces.IAgentRepository
}

func NewAgentService(repo interfaces.IAgentRepository) *AgentService {
	return &AgentService{repo: repo}
}

// CreateBot creates a new bot.
func (s *AgentService) CreateBot(userID uint, botType, targetWallet, chain string, score int, recommendation string, budget, maxPerTrade float64, conditions []string, consensusThreshold, consensusWindowMin, minScore int) (*entities.AgentConfig, *pkg.AppError) {
	condJSON, _ := json.Marshal(conditions)

	config := &entities.AgentConfig{
		UserID:              userID,
		BotType:             botType,
		TargetWalletAddress: targetWallet,
		TargetWalletChain:   chain,
		TargetWalletScore:   score,
		Recommendation:      recommendation,
		Budget:              budget,
		MaxPerTrade:         maxPerTrade,
		Conditions:          condJSON,
		Status:              constants.DefaultAgentStatus,
		ConsensusThreshold:  consensusThreshold,
		ConsensusWindowMin:  consensusWindowMin,
		MinScore:            minScore,
	}

	if err := s.repo.CreateConfig(config); err != nil {
		return nil, pkg.ErrInternal()
	}

	return config, nil
}

// UpdateBot updates a bot's configuration.
func (s *AgentService) UpdateBot(botID, userID uint, budget, maxPerTrade float64, conditions []string, consensusThreshold, consensusWindowMin, minScore int) (*entities.AgentConfig, *pkg.AppError) {
	config, err := s.repo.FindByID(botID)
	if err != nil || config.UserID != userID {
		return nil, pkg.ErrNotFound(constants.DataNotFound)
	}

	if budget > 0 {
		config.Budget = budget
	}
	if maxPerTrade > 0 {
		config.MaxPerTrade = maxPerTrade
	}
	if conditions != nil {
		condJSON, _ := json.Marshal(conditions)
		config.Conditions = condJSON
	}
	if consensusThreshold >= 2 {
		config.ConsensusThreshold = consensusThreshold
	}
	if consensusWindowMin >= 5 {
		config.ConsensusWindowMin = consensusWindowMin
	}
	if minScore >= 0 && minScore <= 100 {
		config.MinScore = minScore
	}

	if err := s.repo.UpdateConfig(config); err != nil {
		return nil, pkg.ErrInternal()
	}

	return config, nil
}

// DeleteBot removes a bot.
func (s *AgentService) DeleteBot(botID, userID uint) *pkg.AppError {
	config, err := s.repo.FindByID(botID)
	if err != nil || config.UserID != userID {
		return pkg.ErrNotFound(constants.DataNotFound)
	}

	if err := s.repo.DeleteConfig(botID); err != nil {
		return pkg.ErrInternal()
	}

	return nil
}

// StartBot activates a bot.
func (s *AgentService) StartBot(botID, userID uint) (*entities.AgentConfig, *pkg.AppError) {
	config, err := s.repo.FindByID(botID)
	if err != nil || config.UserID != userID {
		return nil, pkg.ErrNotFound(constants.DataNotFound)
	}

	config.Status = "active"
	if err := s.repo.UpdateConfig(config); err != nil {
		return nil, pkg.ErrInternal()
	}

	return config, nil
}

// PauseBot pauses a bot.
func (s *AgentService) PauseBot(botID, userID uint) (*entities.AgentConfig, *pkg.AppError) {
	config, err := s.repo.FindByID(botID)
	if err != nil || config.UserID != userID {
		return nil, pkg.ErrNotFound(constants.DataNotFound)
	}

	config.Status = "paused"
	if err := s.repo.UpdateConfig(config); err != nil {
		return nil, pkg.ErrInternal()
	}

	return config, nil
}

// GetBot returns a single bot.
func (s *AgentService) GetBot(botID, userID uint) (*entities.AgentConfig, *pkg.AppError) {
	config, err := s.repo.FindByID(botID)
	if err != nil || config.UserID != userID {
		return nil, pkg.ErrNotFound(constants.DataNotFound)
	}
	return config, nil
}

// ListBots returns all bots for a user.
func (s *AgentService) ListBots(userID uint) ([]entities.AgentConfig, *pkg.AppError) {
	configs, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, pkg.ErrInternal()
	}
	return configs, nil
}

// GetTrades returns trade history for a bot.
func (s *AgentService) GetTrades(botID, userID uint, limit int) ([]entities.AgentTrade, *pkg.AppError) {
	config, err := s.repo.FindByID(botID)
	if err != nil || config.UserID != userID {
		return nil, pkg.ErrNotFound(constants.DataNotFound)
	}

	trades, err := s.repo.FindTradesByConfigID(botID, limit)
	if err != nil {
		return nil, pkg.ErrInternal()
	}

	return trades, nil
}
