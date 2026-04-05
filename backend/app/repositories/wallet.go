package repositories

import (
	"miora-ai/app/entities"

	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) FindByAddress(address string) (*entities.Wallet, error) {
	var wallet entities.Wallet
	if err := r.db.Where("address = ?", address).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) Create(wallet *entities.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *WalletRepository) SaveTransactions(txs []entities.Transaction) error {
	if len(txs) == 0 {
		return nil
	}
	return r.db.Save(&txs).Error
}

func (r *WalletRepository) GetTransactions(walletID uint) ([]entities.Transaction, error) {
	var txs []entities.Transaction
	if err := r.db.Where("wallet_id = ?", walletID).Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *WalletRepository) SaveMetric(metric *entities.WalletMetric) error {
	return r.db.Save(metric).Error
}

func (r *WalletRepository) GetMetric(walletID uint) (*entities.WalletMetric, error) {
	var metric entities.WalletMetric
	if err := r.db.Where("wallet_id = ?", walletID).First(&metric).Error; err != nil {
		return nil, err
	}
	return &metric, nil
}
