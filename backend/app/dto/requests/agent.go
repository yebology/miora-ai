package requests

// CreateBot is the request body for POST /agent/bots.
// BotType determines which fields are required:
//   - "wallet": target_wallet_address, target_wallet_chain required
//   - "consensus": consensus_threshold, consensus_window_min, min_score required
type CreateBot struct {
	BotType             string   `json:"bot_type" validate:"required,oneof=wallet consensus"`
	TargetWalletAddress string   `json:"target_wallet_address"`
	TargetWalletChain   string   `json:"target_wallet_chain"`
	TargetWalletScore   int      `json:"target_wallet_score"`
	Recommendation      string   `json:"recommendation"`
	Budget              float64  `json:"budget" validate:"required,gt=0"`
	MaxPerTrade         float64  `json:"max_per_trade" validate:"required,gt=0"`
	Conditions          []string `json:"conditions"`
	ConsensusThreshold  int      `json:"consensus_threshold" validate:"omitempty,gte=2,lte=20"`
	ConsensusWindowMin  int      `json:"consensus_window_min" validate:"omitempty,gte=5,lte=1440"`
	MinScore            int      `json:"min_score" validate:"omitempty,gte=0,lte=100"`
}

// UpdateBot is the request body for PUT /agent/bots/:id.
type UpdateBot struct {
	Budget             float64  `json:"budget" validate:"omitempty,gte=0"`
	MaxPerTrade        float64  `json:"max_per_trade" validate:"omitempty,gte=0"`
	Conditions         []string `json:"conditions"`
	ConsensusThreshold int      `json:"consensus_threshold" validate:"omitempty,gte=2,lte=20"`
	ConsensusWindowMin int      `json:"consensus_window_min" validate:"omitempty,gte=5,lte=1440"`
	MinScore           int      `json:"min_score" validate:"omitempty,gte=0,lte=100"`
}
