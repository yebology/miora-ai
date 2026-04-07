package interfaces

import "miora-ai/app/entities"

type IWatchlistRepository interface {
	Create(item *entities.Watchlist) error
	Delete(userID uint, walletAddress string) error
	Update(userID uint, walletAddress string, updates map[string]interface{}) error
	FindByUser(userID uint) ([]entities.Watchlist, error)
	FindByWallet(walletAddress string) ([]entities.Watchlist, error)
	Exists(userID uint, walletAddress string) (bool, error)
}
