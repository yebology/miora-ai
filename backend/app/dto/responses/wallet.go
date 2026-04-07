package responses

import "time"

type WalletAnalysis struct {
	Address           string        `json:"address"`
	Chain             string        `json:"chain"`
	TotalTransactions int           `json:"total_transactions"`
	ProfitConsistency float64       `json:"profit_consistency"`
	WinRate           float64       `json:"win_rate"`
	RiskExposure      float64       `json:"risk_exposure"`
	EntryTiming       float64       `json:"entry_timing"`
	TokenQuality      float64       `json:"token_quality"`
	TradeDiscipline   float64       `json:"trade_discipline"`
	FinalScore        float64       `json:"final_score"`
	Recommendation    string        `json:"recommendation"`
	AiInsight         string        `json:"ai_insight,omitempty"`
	TradedTokens      []TradedToken `json:"traded_tokens,omitempty"`
	Conditions        []Condition   `json:"conditions,omitempty"`
}

// TradedToken represents a token the wallet traded, with PnL data.
type TradedToken struct {
	ContractAddress string     `json:"contract_address"`
	Symbol          string     `json:"symbol"`
	Chain           string     `json:"chain"`
	PnlPercent      float64    `json:"pnl_percent"`
	BuyPrice        float64    `json:"buy_price"`
	ExitPrice       float64    `json:"exit_price"`
	BuyTime         time.Time  `json:"buy_time"`
	ExitTime        *time.Time `json:"exit_time,omitempty"` // nil if unrealized (still holding)
	Status          string     `json:"status"`
}

// Condition represents a suggested filter for conditional follow notifications.
// User can select which conditions to apply — only trades matching all selected
// conditions will trigger a notification.
type Condition struct {
	ID       string      `json:"id"`       // Unique identifier (e.g. "min_liquidity")
	Label    string      `json:"label"`    // Human-readable description
	Type     string      `json:"type"`     // Data type: "number"
	Field    string      `json:"field"`    // DexScreener field to check (e.g. "liquidity", "pair_age_hours", "market_cap")
	Operator string      `json:"operator"` // Comparison: "gte" (>=), "lte" (<=)
	Value    interface{} `json:"value"`    // Threshold value
}
