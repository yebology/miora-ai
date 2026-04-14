package dto

// SwapQuote holds the result of a DEX quote request.
type SwapQuote struct {
	InputMint   string `json:"input_mint"`
	OutputMint  string `json:"output_mint"`
	InAmount    string `json:"in_amount"`
	OutAmount   string `json:"out_amount"`
	PriceImpact string `json:"price_impact,omitempty"`
	SlippageBps int    `json:"slippage_bps"`
	Chain       string `json:"chain"`
	Dex         string `json:"dex"` // "1inch"
}
