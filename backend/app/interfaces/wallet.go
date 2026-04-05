package interfaces

import "miora-ai/app/entities"

type WalletRepository interface {
	FindByAddress(address string) (*entities.Wallet, error)
	Create(wallet *entities.Wallet) error
	SaveTransactions(txs []entities.Transaction) error
	GetTransactions(walletID uint) ([]entities.Transaction, error)
	SaveMetric(metric *entities.WalletMetric) error
	GetMetric(walletID uint) (*entities.WalletMetric, error)
}
