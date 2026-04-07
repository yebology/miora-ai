package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"miora-ai/app/interfaces"
)

// AlchemySolana fetches transaction signatures from Solana mainnet
// using Alchemy's getSignaturesForAddress JSON-RPC method.
//
// API docs: https://docs.alchemy.com/reference/getsignaturesforaddress
type AlchemySolana struct {
	apiKey string
}

// NewAlchemySolana creates a new Solana client with the given Alchemy API key.
func NewAlchemySolana(apiKey string) *AlchemySolana {

	return &AlchemySolana{apiKey: apiKey}

}

// solanaTx represents a single transaction signature from the Solana RPC response.
type solanaTx struct {
	Signature string `json:"signature"` // Unique transaction hash (base58 encoded)
	BlockTime int64  `json:"blockTime"` // Unix timestamp when the transaction was processed
	Slot      uint64 `json:"slot"`      // Slot number (Solana's equivalent of block number, ~400ms per slot)
}

// solanaRPCResponse is the JSON-RPC response from getSignaturesForAddress.
type solanaRPCResponse struct {
	Result []solanaTx `json:"result"` // List of transaction signatures
}

// GetTransfers fetches the last 100 transaction signatures for the given address.
//
// Uses getSignaturesForAddress with limit: 100.
//
// Note: This method only returns signatures and timestamps, not full transfer
// details (from, to, value, token). To get complete transfer data, each
// signature would need a follow-up getTransaction call.
//
// Returns a normalized slice of TransferData with only Hash and Timestamp populated.
func (a *AlchemySolana) GetTransfers(address string, chain ...string) ([]interfaces.TransferData, error) {

	url := fmt.Sprintf("https://solana-mainnet.g.alchemy.com/v2/%s", a.apiKey)

	// JSON-RPC request payload
	// params[0]: wallet address to query
	// params[1]: options object with limit (max 100 signatures per request)
	payload := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "getSignaturesForAddress",
		"params":  []interface{}{address, map[string]interface{}{"limit": 100}},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("alchemy solana request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result solanaRPCResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	// Normalize into chain-agnostic TransferData
	// Only Hash and Timestamp are available from this endpoint
	transfers := make([]interfaces.TransferData, 0, len(result.Result))
	for _, tx := range result.Result {
		transfers = append(transfers, interfaces.TransferData{
			Hash:      tx.Signature,
			Timestamp: tx.BlockTime,
		})
	}

	return transfers, nil

}
