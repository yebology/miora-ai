package interfaces

import (
	"miora-ai/app/entities"
	"miora-ai/pkg"
)

// IUserService defines the business logic contract for user operations.
type IUserService interface {
	FindOrCreateByWallet(walletAddress string) (*entities.User, *pkg.AppError)
	UpdateEmail(walletAddress, email string) (*entities.User, *pkg.AppError)
}

// IUserRepository defines the data access contract for user operations.
type IUserRepository interface {
	FindByWalletAddress(walletAddress string) (*entities.User, error)
	FindByID(id uint) (*entities.User, error)
	Create(user *entities.User) error
	Update(user *entities.User) error
}
