// Package services contains user business logic.
//
// Users are identified by wallet address. No Firebase — just connect wallet.
package services

import (
	"errors"

	"miora-ai/app/entities"
	"miora-ai/app/interfaces"
	"miora-ai/constants"
	"miora-ai/pkg"

	"gorm.io/gorm"
)

type UserService struct {
	repo interfaces.IUserRepository
}

func NewUserService(repo interfaces.IUserRepository) *UserService {
	return &UserService{repo: repo}
}

// FindOrCreateByWallet finds a user by wallet address, or creates one if not found.
func (s *UserService) FindOrCreateByWallet(walletAddress string) (*entities.User, *pkg.AppError) {
	user, err := s.repo.FindByWalletAddress(walletAddress)
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, pkg.ErrInternal()
	}

	user = &entities.User{
		WalletAddress: walletAddress,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, pkg.ErrInternal()
	}

	return user, nil
}

// UpdateEmail sets the email for a user (optional, for trade notifications).
func (s *UserService) UpdateEmail(walletAddress, email string) (*entities.User, *pkg.AppError) {
	user, appErr := s.FindOrCreateByWallet(walletAddress)
	if appErr != nil {
		return nil, appErr
	}

	user.Email = email
	if err := s.repo.Update(user); err != nil {
		return nil, pkg.ErrInternal()
	}

	return user, nil
}

// GetByWalletAddress returns a user by wallet address.
func (s *UserService) GetByWalletAddress(walletAddress string) (*entities.User, *pkg.AppError) {
	user, err := s.repo.FindByWalletAddress(walletAddress)
	if err != nil {
		return nil, pkg.ErrNotFound(constants.DataNotFound)
	}
	return user, nil
}
