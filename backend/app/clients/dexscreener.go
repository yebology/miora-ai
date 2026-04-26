// DexScreener client for fetching token pair data (price, liquidity, volume, market cap).
// Free API, no key required. Rate limit: 300 req/min for pair/token endpoints.
//
// API docs: https://docs.dexscreener.com/api/reference
package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"miora-ai/app/dto"
	"net/http"
)

// DexScreener fetches token pair data from the DexScreener API.
type DexScreener struct{}

// NewDexScreener creates a new DexScreener client.
func NewDexScreener() *DexScreener {

	return &DexScreener{}

}

// dexPair represents a single trading pair from the DexScreener API response.
type dexPair struct {
	ChainID       string   `json:"chainId"`       // Blockchain network (e.g. "ethereum", "solana")
	PairAddress   string   `json:"pairAddress"`   // Address of the liquidity pool on the DEX
	PriceUsd      string   `json:"priceUsd"`      // Current token price in USD
	Liquidity     dexLiq   `json:"liquidity"`     // Total money locked in the pool — higher = safer to trade, lower = risky/easy to manipulate
	Volume        dexVol   `json:"volume"`        // Total USD value of trades in a time window — shows how actively the token is being traded
	FDV           float64  `json:"fdv"`           // Fully Diluted Valuation — market cap if ALL tokens (including locked/unvested) were in circulation
	MarketCap     float64  `json:"marketCap"`     // Current market cap — price × circulating supply. Shows how "big" the token is
	PairCreatedAt int64    `json:"pairCreatedAt"` // Unix timestamp (ms) when this trading pair was first created on the DEX
	PriceChange   dexPChg  `json:"priceChange"`   // Price change percentages across time windows (1h, 6h, 24h)
	BaseToken     dexToken `json:"baseToken"`     // The token being traded (not the quote token like USDC/WETH)
}

type dexLiq struct {
	USD float64 `json:"usd"`
}

type dexVol struct {
	H24 float64 `json:"h24"`
}

type dexPChg struct {
	H24 float64 `json:"h24"`
}

type dexToken struct {
	Address string `json:"address"`
	Symbol  string `json:"symbol"`
}

type dexWrappedResponse struct {
	Pairs []dexPair `json:"pairs"`
}

// GetTokenPairs fetches all trading pairs for a token on a given chain.
// Returns normalized TokenPairData for chain-agnostic processing.
func (d *DexScreener) GetTokenPairs(chainID, tokenAddress string) ([]dto.TokenPairData, error) {

	url := fmt.Sprintf("https://api.dexscreener.com/token-pairs/v1/%s/%s", chainID, tokenAddress)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("dexscreener request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var pairs []dexPair
	if err := json.Unmarshal(body, &pairs); err != nil {
		var wrapped dexWrappedResponse
		if err2 := json.Unmarshal(body, &wrapped); err2 != nil {
			return nil, fmt.Errorf("unmarshal response: %w", err)
		}
		pairs = wrapped.Pairs
	}

	// Normalize to interface type
	result := make([]dto.TokenPairData, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, dto.TokenPairData{
			ChainID:        p.ChainID,
			PairAddress:    p.PairAddress,
			PriceUsd:       p.PriceUsd,
			Liquidity:      p.Liquidity.USD,
			VolumeH24:      p.Volume.H24,
			FDV:            p.FDV,
			MarketCap:      p.MarketCap,
			PairCreatedAt:  p.PairCreatedAt,
			PriceChangeH24: p.PriceChange.H24,
			BaseSymbol:     p.BaseToken.Symbol,
			BaseAddress:    p.BaseToken.Address,
		})
	}

	return result, nil

}
