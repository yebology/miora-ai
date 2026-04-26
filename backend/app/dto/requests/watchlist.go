package requests

type FollowWallet struct {
	WalletAddress  string   `json:"wallet_address" validate:"required"`
	Chain          string   `json:"chain" validate:"required,oneof=evm ethereum arbitrum optimism base polygon"`
	Recommendation string   `json:"recommendation"`
	Conditions     []string `json:"conditions"` // Selected condition IDs (e.g. ["min_liquidity", "min_pair_age"])
	EmailNotify    bool     `json:"email_notify"`
}

// UpdateWatchlist is the request body for PUT /watchlist/:address.
type UpdateWatchlist struct {
	Conditions  []string `json:"conditions"`
	EmailNotify *bool    `json:"email_notify"` // Pointer so we can distinguish "not sent" from "false"
}
