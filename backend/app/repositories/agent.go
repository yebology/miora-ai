package repositories

import (
	"miora-ai/app/entities"

	"gorm.io/gorm"
)

type AgentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) *AgentRepository {
	return &AgentRepository{db: db}
}

func (r *AgentRepository) FindByID(id uint) (*entities.AgentConfig, error) {
	var config entities.AgentConfig
	if err := r.db.First(&config, id).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *AgentRepository) FindByUserID(userID uint) ([]entities.AgentConfig, error) {
	var configs []entities.AgentConfig
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *AgentRepository) FindActiveConfigs() ([]entities.AgentConfig, error) {
	var configs []entities.AgentConfig
	if err := r.db.Where("status = ?", "active").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *AgentRepository) CreateConfig(config *entities.AgentConfig) error {
	return r.db.Create(config).Error
}

func (r *AgentRepository) UpdateConfig(config *entities.AgentConfig) error {
	return r.db.Save(config).Error
}

func (r *AgentRepository) DeleteConfig(id uint) error {
	return r.db.Delete(&entities.AgentConfig{}, id).Error
}

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

func (r *AgentRepository) CreateTrade(trade *entities.AgentTrade) error {
	return r.db.Create(trade).Error
}
