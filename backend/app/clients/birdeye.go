// Birdeye client for fetching Solana token historical prices.
// Requires API key (free tier: ~100 requests/day).
// Used for PnL calculation on Solana where Moralis doesn't support historical prices.
//
// API docs: https://docs.birdeye.so/docs/premium-apis-1
package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"miora-ai/app/dto"
)

// Birdeye fetches historical token prices from the Birdeye API.
type Birdeye struct {
	apiKey string
}

// NewBirdeye creates a new Birdeye client with the given API key.
func NewBirdeye(apiKey string) *Birdeye {

	return &Birdeye{apiKey: apiKey}

}

type birdeyePriceResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Value float64 `json:"value"`
	} `json:"data"`
}

// GetHistoricalPrice fetches the token price closest to the given unix timestamp.
//
// Endpoint: GET https://public-api.birdeye.so/defi/historical_price_unix
// Params:
//   - address: Solana token mint address
//   - address_type: "token"
//   - time_from: unix timestamp to query
//
// Returns the price at the closest available time to the given timestamp.
func (b *Birdeye) GetHistoricalPrice(tokenAddress string, unixTimestamp int64) (*dto.TokenPriceData, error) {

	url := fmt.Sprintf(
		"https://public-api.birdeye.so/defi/historical_price_unix?address=%s&address_type=token&time_from=%d",
		tokenAddress, unixTimestamp,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("X-API-KEY", b.apiKey)
	req.Header.Set("x-chain", "solana")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("birdeye request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("birdeye error %d: %s", resp.StatusCode, string(body))
	}

	var result birdeyePriceResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("birdeye: unsuccessful response")
	}

	return &dto.TokenPriceData{
		TokenAddress: tokenAddress,
		UsdPrice:     result.Data.Value,
	}, nil

}
