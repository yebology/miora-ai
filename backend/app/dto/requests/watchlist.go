package requests

type FollowWallet struct {
	WalletAddress  string   `json:"wallet_address" validate:"required"`
	Chain          string   `json:"chain" validate:"required,oneof=evm svm ethereum arbitrum optimism base polygon solana"`
	Recommendation string   `json:"recommendation"`
	Conditions     []string `json:"conditions"` // Selected condition IDs (e.g. ["min_liquidity", "min_pair_age"])
	EmailNotify    bool     `json:"email_notify"`
}
