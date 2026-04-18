// Package clients provides the AgentKit sidecar HTTP client.
//
// This client communicates with the Python AgentKit sidecar service
// running on localhost:8090 to execute autonomous trading operations
// via Coinbase AgentKit on Base Sepolia.
package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AgentKitClient communicates with the Python AgentKit sidecar.
type AgentKitClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAgentKitClient creates a new AgentKit sidecar client.
// sidecarURL is the base URL of the Python sidecar (e.g. "http://localhost:8090").
func NewAgentKitClient(sidecarURL string) *AgentKitClient {
	if sidecarURL == "" {
		sidecarURL = "http://localhost:8090"
	}
	return &AgentKitClient{
		baseURL: sidecarURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AgentWalletInfo holds the agent wallet details from the sidecar.
type AgentWalletInfo struct {
	Address string `json:"address"`
	Network string `json:"network"`
	Balance string `json:"balance"`
}

// GetWallet returns the agent's wallet address and balance.
func (c *AgentKitClient) GetWallet() (*AgentWalletInfo, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/wallet")
	if err != nil {
		return nil, fmt.Errorf("agent sidecar unreachable: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("agent sidecar error: %s", string(body))
	}

	var info AgentWalletInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse wallet info: %w", err)
	}

	return &info, nil
}

// IsHealthy checks if the agent sidecar is running and ready.
func (c *AgentKitClient) IsHealthy() bool {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// SwapRequest is the request body for the sidecar swap endpoint.
type SwapRequest struct {
	TokenAddress string `json:"token_address"`
	AmountETH    string `json:"amount_eth"`
	TokenSymbol  string `json:"token_symbol"`
	Direction    string `json:"direction"` // "buy" or "sell"
}

// SwapResult is the response from the sidecar swap endpoint.
type SwapResult struct {
	Status       string `json:"status"`
	TokenAddress string `json:"token_address"`
	TokenSymbol  string `json:"token_symbol"`
	AmountETH    string `json:"amount_eth"`
	Result       string `json:"result"`
	AgentWallet  string `json:"agent_wallet"`
}

// ExecuteSwap calls the Python sidecar to execute a token swap via AgentKit.
func (c *AgentKitClient) ExecuteSwap(tokenAddress, tokenSymbol, amountETH, direction string) (*SwapResult, error) {
	payload, _ := json.Marshal(SwapRequest{
		TokenAddress: tokenAddress,
		AmountETH:    amountETH,
		TokenSymbol:  tokenSymbol,
		Direction:    direction,
	})

	resp, err := c.httpClient.Post(c.baseURL+"/swap", "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("agent sidecar unreachable: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("swap failed: %s", string(body))
	}

	var result SwapResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse swap result: %w", err)
	}

	return &result, nil
}
