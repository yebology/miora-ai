// Package repositories provides data access implementations using GORM.
//
// Each repository implements an interface from app/interfaces/.
// Repository methods return Go's standard error — the service layer
// is responsible for converting these into *pkg.AppError.
package repositories

import (
	"miora-ai/app/entities"

	"gorm.io/gorm"
)

// WalletRepository implements interfaces.IWalletRepository.
// Handles all database operations for wallets, transactions, and metrics.
type WalletRepository struct {
	db *gorm.DB
}

// NewWalletRepository creates a new WalletRepository with the given database connection.
func NewWalletRepository(db *gorm.DB) *WalletRepository {

	return &WalletRepository{db: db}

}

// FindByAddress looks up a wallet by its on-chain address.
// Returns gorm.ErrRecordNotFound if the wallet doesn't exist.
func (r *WalletRepository) FindByAddress(address string) (*entities.Wallet, error) {

	var wallet entities.Wallet
	if err := r.db.Where("address = ?", address).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil

}

// Create inserts a new wallet record into the database.
// The wallet's ID is populated after creation.
func (r *WalletRepository) Create(wallet *entities.Wallet) error {

	return r.db.Create(wallet).Error

}

// SaveTransactions upserts a batch of transactions.
// Skips if the slice is empty.
func (r *WalletRepository) SaveTransactions(txs []entities.Transaction) error {

	if len(txs) == 0 {
		return nil
	}
	return r.db.Save(&txs).Error

}

// GetTransactions returns all transactions for a given wallet ID.
func (r *WalletRepository) GetTransactions(walletID uint) ([]entities.Transaction, error) {

	var txs []entities.Transaction
	if err := r.db.Where("wallet_id = ?", walletID).Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil

}

// SaveMetric upserts the wallet metric record (scoring data).
func (r *WalletRepository) SaveMetric(metric *entities.WalletMetric) error {

	return r.db.Save(metric).Error

}

// GetMetric returns the scoring metric for a given wallet ID.
// Returns gorm.ErrRecordNotFound if no metric exists yet.
func (r *WalletRepository) GetMetric(walletID uint) (*entities.WalletMetric, error) {

	var metric entities.WalletMetric
	if err := r.db.Where("wallet_id = ?", walletID).First(&metric).Error; err != nil {
		return nil, err
	}
	return &metric, nil

}

// FindAllWithMetrics returns all wallets that have a metric with score >= minScore.
// Used by the consensus agent to scan all analyzed wallets in the database.
func (r *WalletRepository) FindAllWithMetrics(minScore int) ([]entities.Wallet, error) {
	var wallets []entities.Wallet
	err := r.db.
		Joins("JOIN wallet_metrics ON wallet_metrics.wallet_id = wallets.id").
		Where("wallet_metrics.final_score >= ?", minScore).
		Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

// Delete removes a wallet and its associated transactions by wallet ID.
func (r *WalletRepository) Delete(walletID uint) error {
	r.db.Where("wallet_id = ?", walletID).Delete(&entities.Transaction{})
	return r.db.Delete(&entities.Wallet{}, walletID).Error
}

// DeleteTransactions removes all transactions for a wallet ID.
func (r *WalletRepository) DeleteTransactions(walletID uint) error {
	return r.db.Where("wallet_id = ?", walletID).Delete(&entities.Transaction{}).Error
}

// DeleteMetric removes the wallet metric for a wallet ID.
func (r *WalletRepository) DeleteMetric(walletID uint) error {
	return r.db.Where("wallet_id = ?", walletID).Delete(&entities.WalletMetric{}).Error
}
