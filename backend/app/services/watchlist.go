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

// UpdateConditions updates the selected conditions and/or email notify preference for a watchlist item.
func (s *WatchlistService) UpdateConditions(userID uint, walletAddress string, conditions []string, emailNotify *bool) *pkg.AppError {

	exists, err := s.repo.Exists(userID, walletAddress)
	if err != nil {
		return pkg.ErrInternal()
	}
	if !exists {
		return pkg.ErrNotFound("Not following this wallet.")
	}

	updates := make(map[string]interface{})

	if conditions != nil {
		condJSON, _ := json.Marshal(conditions)
		updates["conditions"] = condJSON
	}

	if emailNotify != nil {
		updates["email_notify"] = *emailNotify
	}

	if len(updates) == 0 {
		return nil
	}

	if err := s.repo.Update(userID, walletAddress, updates); err != nil {
		return pkg.ErrInternal()
	}

	return nil

}
