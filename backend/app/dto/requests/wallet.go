package requests

type AnalyzeWallet struct {
	Address string `json:"address" validate:"required"`
	Chain   string `json:"chain" validate:"required,oneof=evm svm ethereum arbitrum optimism base polygon solana"`
	Limit   int    `json:"limit"` // Max transactions to fetch (default 25, max 50)
}
