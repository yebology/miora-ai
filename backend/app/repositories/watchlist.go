package repositories

import (
	"miora-ai/app/entities"

	"gorm.io/gorm"
)

type WatchlistRepository struct {
	db *gorm.DB
}

func NewWatchlistRepository(db *gorm.DB) *WatchlistRepository {

	return &WatchlistRepository{db: db}

}

func (r *WatchlistRepository) Create(item *entities.Watchlist) error {

	return r.db.Create(item).Error

}

func (r *WatchlistRepository) Delete(userID uint, walletAddress string) error {

	return r.db.Where("user_id = ? AND wallet_address = ?", userID, walletAddress).Delete(&entities.Watchlist{}).Error

}

func (r *WatchlistRepository) FindByUser(userID uint) ([]entities.Watchlist, error) {

	var items []entities.Watchlist
	if err := r.db.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil

}

func (r *WatchlistRepository) FindByWallet(walletAddress string) ([]entities.Watchlist, error) {

	var items []entities.Watchlist
	query := r.db
	if walletAddress != "" {
		query = query.Where("wallet_address = ?", walletAddress)
	}
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil

}

func (r *WatchlistRepository) Exists(userID uint, walletAddress string) (bool, error) {

	var count int64
	err := r.db.Model(&entities.Watchlist{}).Where("user_id = ? AND wallet_address = ?", userID, walletAddress).Count(&count).Error
	return count > 0, err

}
