package responses

type WalletAnalysis struct {
	Address           string  `json:"address"`
	Chain             string  `json:"chain"`
	TotalTransactions int     `json:"total_transactions"`
	ProfitConsistency float64 `json:"profit_consistency"`
	WinRate           float64 `json:"win_rate"`
	RiskExposure      float64 `json:"risk_exposure"`
	EntryTiming       float64 `json:"entry_timing"`
	TokenQuality      float64 `json:"token_quality"`
	TradeDiscipline   float64 `json:"trade_discipline"`
	FinalScore        float64 `json:"final_score"`
	Recommendation    string  `json:"recommendation"`
	AiInsight         string  `json:"ai_insight,omitempty"`
}
