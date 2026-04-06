package interfaces

import (
	"miora-ai/app/dto/responses"
	"miora-ai/app/entities"
	"miora-ai/pkg"
)

// IWalletService defines the business logic contract for wallet operations.
// Methods return *pkg.AppError instead of Go's error for structured error handling.
type IWalletService interface {
	AnalyzeWallet(address, chain string) (*responses.WalletAnalysis, *pkg.AppError)
}

// IWalletRepository defines the data access contract for wallet operations.
// Methods return Go's error — the service layer converts these to *pkg.AppError.
type IWalletRepository interface {
	FindByAddress(address string) (*entities.Wallet, error)
	Create(wallet *entities.Wallet) error
	SaveTransactions(txs []entities.Transaction) error
	GetTransactions(walletID uint) ([]entities.Transaction, error)
	SaveMetric(metric *entities.WalletMetric) error
	GetMetric(walletID uint) (*entities.WalletMetric, error)
}
