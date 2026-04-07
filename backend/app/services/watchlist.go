package services

import (
	"encoding/json"

	"miora-ai/app/entities"
	"miora-ai/app/interfaces"
	"miora-ai/constants"
	"miora-ai/pkg"
)

type WatchlistService struct {
	repo interfaces.IWatchlistRepository
}

func NewWatchlistService(repo interfaces.IWatchlistRepository) *WatchlistService {

	return &WatchlistService{repo: repo}

}

// Follow adds a wallet to the user's watchlist.
func (s *WatchlistService) Follow(userID uint, walletAddress, chain, recommendation string, conditions []string, emailNotify bool) *pkg.AppError {

	exists, err := s.repo.Exists(userID, walletAddress)
	if err != nil {
		return pkg.ErrInternal()
	}
	if exists {
		return pkg.ErrConflict("Already following this wallet.")
	}

	condJSON, _ := json.Marshal(conditions)

	item := &entities.Watchlist{
		UserID:         userID,
		WalletAddress:  walletAddress,
		Chain:          chain,
		Recommendation: recommendation,
		Conditions:     condJSON,
		EmailNotify:    emailNotify,
	}

	if err := s.repo.Create(item); err != nil {
		return pkg.ErrInternal()
	}

	return nil

}

// Unfollow removes a wallet from the user's watchlist.
func (s *WatchlistService) Unfollow(userID uint, walletAddress string) *pkg.AppError {

	if err := s.repo.Delete(userID, walletAddress); err != nil {
		return pkg.ErrInternal()
	}
	return nil

}

// GetUserWatchlist returns all wallets the user is following.
func (s *WatchlistService) GetUserWatchlist(userID uint) ([]entities.Watchlist, *pkg.AppError) {

	items, err := s.repo.FindByUser(userID)
	if err != nil {
		return nil, pkg.ErrNotFound(constants.DataNotFound)
	}
	return items, nil

}
