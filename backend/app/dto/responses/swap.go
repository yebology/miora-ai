package responses

// SwapQuote is the response from POST /api/swap/quote.
//
// Example response:
//
//	{
//	  "chain": "svm",
//	  "input_mint": "So11111111111111111111111111111111111111112",
//	  "output_mint": "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
//	  "input_amount": "100000000",
//	  "output_amount": "14523000",
//	  "price_impact": "0.01",
//	  "route": "Raydium → Orca"
//	}
type SwapQuote struct {
	Chain        string `json:"chain"`                  // Blockchain used: "evm" or "svm"
	InputMint    string `json:"input_mint"`             // Token address being sold
	OutputMint   string `json:"output_mint"`            // Token address being bought
	InputAmount  string `json:"input_amount"`           // Amount being sold (smallest unit)
	OutputAmount string `json:"output_amount"`          // Amount received (smallest unit)
	PriceImpact  string `json:"price_impact,omitempty"` // How much the swap moves the market price (e.g. "0.01" = 0.01%). Higher = worse.
	Route        string `json:"route,omitempty"`        // DEX route taken (e.g. "Raydium → Orca"). Shows which DEXes the aggregator used.
}
