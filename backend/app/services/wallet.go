package services

import (
	"errors"
	"time"

	"miora-ai/app/dto/responses"
	"miora-ai/app/entities"
	"miora-ai/app/interfaces"

	"gorm.io/gorm"
)

type WalletService struct {
	repo      interfaces.WalletRepository
	evmClient interfaces.BlockchainClient
	svmClient interfaces.BlockchainClient
}

func NewWalletService(
	repo interfaces.WalletRepository,
	evmClient interfaces.BlockchainClient,
	svmClient interfaces.BlockchainClient,
) *WalletService {
	return &WalletService{
		repo:      repo,
		evmClient: evmClient,
		svmClient: svmClient,
	}
}

func (s *WalletService) AnalyzeWallet(address, chain string) (*responses.WalletAnalysis, error) {
	// Find or create wallet
	wallet, err := s.repo.FindByAddress(address)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		wallet = &entities.Wallet{Address: address, Chain: chain}
		if err := s.repo.Create(wallet); err != nil {
			return nil, err
		}
	}

	// Pick the right client
	var client interfaces.BlockchainClient
	switch chain {
	case "evm":
		client = s.evmClient
	case "svm":
		client = s.svmClient
	default:
		return nil, errors.New("unsupported chain")
	}

	// Fetch transfers from Alchemy
	transfers, err := client.GetTransfers(address)
	if err != nil {
		return nil, err
	}

	// Convert to entities and save
	txEntities := make([]entities.Transaction, 0, len(transfers))
	for _, t := range transfers {
		txEntities = append(txEntities, entities.Transaction{
			WalletID:    wallet.ID,
			Hash:        t.Hash,
			Chain:       chain,
			From:        t.From,
			To:          t.To,
			Value:       t.Value,
			TokenSymbol: t.TokenSymbol,
			BlockNumber: t.BlockNumber,
			Timestamp:   time.Unix(t.Timestamp, 0),
		})
	}

	if err := s.repo.SaveTransactions(txEntities); err != nil {
		return nil, err
	}

	// Calculate metrics
	metric := s.calculateMetrics(wallet.ID, txEntities)
	if err := s.repo.SaveMetric(metric); err != nil {
		return nil, err
	}

	return &responses.WalletAnalysis{
		Address:           address,
		Chain:             chain,
		TotalTransactions: metric.TotalTransactions,
		ProfitConsistency: metric.ProfitConsistency,
		WinRate:           metric.WinRate,
		RiskExposure:      metric.RiskExposure,
		EntryTiming:       metric.EntryTiming,
		TokenQuality:      metric.TokenQuality,
		TradeDiscipline:   metric.TradeDiscipline,
		FinalScore:        metric.FinalScore,
		Recommendation:    metric.Recommendation,
	}, nil
}
