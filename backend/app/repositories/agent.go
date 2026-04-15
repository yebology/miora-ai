// Package repositories provides the data access layer for agent operations.
package repositories

import (
	"miora-ai/app/entities"

	"gorm.io/gorm"
)

// AgentRepository implements interfaces.IAgentRepository.
type AgentRepository struct {
	db *gorm.DB
}

// NewAgentRepository creates a new AgentRepository.
func NewAgentRepository(db *gorm.DB) *AgentRepository {
	return &AgentRepository{db: db}
}

// FindConfigByUserID returns the agent config for a user.
func (r *AgentRepository) FindConfigByUserID(userID uint) (*entities.AgentConfig, error) {
	var config entities.AgentConfig
	if err := r.db.Where("user_id = ?", userID).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// CreateConfig creates a new agent config.
func (r *AgentRepository) CreateConfig(config *entities.AgentConfig) error {
	return r.db.Create(config).Error
}

// UpdateConfig updates an existing agent config.
func (r *AgentRepository) UpdateConfig(config *entities.AgentConfig) error {
	return r.db.Save(config).Error
}

// FindActiveConfigs returns all agent configs with status "active".
func (r *AgentRepository) FindActiveConfigs() ([]entities.AgentConfig, error) {
	var configs []entities.AgentConfig
	if err := r.db.Where("status = ?", "active").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// FindTradesByConfigID returns trades for an agent config, ordered by most recent.
func (r *AgentRepository) FindTradesByConfigID(configID uint, limit int) ([]entities.AgentTrade, error) {
	var trades []entities.AgentTrade
	query := r.db.Where("agent_config_id = ?", configID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&trades).Error; err != nil {
		return nil, err
	}
	return trades, nil
}

// CreateTrade records a new agent trade.
func (r *AgentRepository) CreateTrade(trade *entities.AgentTrade) error {
	return r.db.Create(trade).Error
}
