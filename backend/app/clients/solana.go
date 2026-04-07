package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"miora-ai/app/interfaces"
)

// AlchemySolana fetches transaction data from Solana mainnet.
// Uses getSignaturesForAddress to get tx list, then getTransaction per signature for details.
type AlchemySolana struct {
	apiKey string
}

func NewAlchemySolana(apiKey string) *AlchemySolana {

	return &AlchemySolana{apiKey: apiKey}

}

type solanaTx struct {
	Signature string `json:"signature"`
	BlockTime int64  `json:"blockTime"`
	Slot      uint64 `json:"slot"`
}

type solanaRPCResponse struct {
	Result []solanaTx `json:"result"`
}

// solanaTransactionDetail holds parsed transaction details.
type solanaTransactionDetail struct {
	Meta struct {
		PreTokenBalances  []solanaTokenBalance `json:"preTokenBalances"`
		PostTokenBalances []solanaTokenBalance `json:"postTokenBalances"`
	} `json:"meta"`
	Transaction struct {
		Message struct {
			AccountKeys []string `json:"accountKeys"`
		} `json:"message"`
	} `json:"transaction"`
	BlockTime int64 `json:"blockTime"`
}

type solanaTokenBalance struct {
	AccountIndex  int    `json:"accountIndex"`
	Mint          string `json:"mint"`
	UITokenAmount struct {
		UIAmount float64 `json:"uiAmount"`
	} `json:"uiTokenAmount"`
	Owner string `json:"owner"`
}

type solanaDetailResponse struct {
	Result *solanaTransactionDetail `json:"result"`
}

// GetTransfers fetches signatures then enriches with transaction details.
func (a *AlchemySolana) GetTransfers(address string, limit int, chain ...string) ([]interfaces.TransferData, error) {

	url := fmt.Sprintf("https://solana-mainnet.g.alchemy.com/v2/%s", a.apiKey)

	if limit <= 0 || limit > 50 {
		limit = 25
	}

	// Step 1: Get signatures
	sigPayload := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "getSignaturesForAddress",
		"params":  []interface{}{address, map[string]interface{}{"limit": limit}},
	}

	sigs, err := a.rpcCall(url, sigPayload)
	if err != nil {
		return nil, fmt.Errorf("get signatures: %w", err)
	}

	var sigResult solanaRPCResponse
	if err := json.Unmarshal(sigs, &sigResult); err != nil {
		return nil, fmt.Errorf("unmarshal signatures: %w", err)
	}

	// Step 2: Get transaction details per signature (limited to avoid rate limits)
	transfers := make([]interfaces.TransferData, 0, len(sigResult.Result))
	maxDetailFetch := limit // Fetch details for all signatures

	for i, sig := range sigResult.Result {
		// Only fetch details for first N signatures
		if i < maxDetailFetch {
			// Small delay between calls to respect rate limits
			if i > 0 {
				time.Sleep(50 * time.Millisecond)
			}

			detail, err := a.getTransactionDetail(url, sig.Signature)
			if err == nil && detail != nil {
				parsed := a.parseTokenTransfers(address, sig.Signature, detail)
				if len(parsed) > 0 {
					transfers = append(transfers, parsed...)
					continue
				}
			}
		}

		// Fallback: basic data without details
		transfers = append(transfers, interfaces.TransferData{
			Hash:      sig.Signature,
			Timestamp: sig.BlockTime,
		})
	}

	return transfers, nil

}

// getTransactionDetail fetches full transaction data for a signature.
func (a *AlchemySolana) getTransactionDetail(url, signature string) (*solanaTransactionDetail, error) {

	payload := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "getTransaction",
		"params": []interface{}{
			signature,
			map[string]interface{}{"encoding": "json", "maxSupportedTransactionVersion": 0},
		},
	}

	body, err := a.rpcCall(url, payload)
	if err != nil {
		return nil, err
	}

	var result solanaDetailResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Result, nil

}

// parseTokenTransfers extracts token transfers by comparing pre/post token balances.
// If the wallet's balance increased → "in" (buy). If decreased → "out" (sell).
func (a *AlchemySolana) parseTokenTransfers(walletAddress, signature string, detail *solanaTransactionDetail) []interfaces.TransferData {

	var transfers []interfaces.TransferData

	// Build pre-balance map: mint → amount for this wallet
	preBalances := make(map[string]float64)
	for _, b := range detail.Meta.PreTokenBalances {
		if b.Owner == walletAddress {
			preBalances[b.Mint] = b.UITokenAmount.UIAmount
		}
	}

	// Compare with post-balances to detect transfers
	for _, b := range detail.Meta.PostTokenBalances {
		if b.Owner != walletAddress {
			continue
		}

		pre := preBalances[b.Mint]
		post := b.UITokenAmount.UIAmount
		diff := post - pre

		if diff == 0 {
			continue
		}

		direction := "in"
		value := diff
		if diff < 0 {
			direction = "out"
			value = -diff
		}

		transfers = append(transfers, interfaces.TransferData{
			Hash:            signature,
			ContractAddress: b.Mint,
			Value:           fmt.Sprintf("%f", value),
			Direction:       direction,
			Timestamp:       detail.BlockTime,
			From:            walletAddress,
			To:              walletAddress,
		})
	}

	return transfers

}

// rpcCall makes a JSON-RPC call and returns the raw response body.
func (a *AlchemySolana) rpcCall(url string, payload interface{}) ([]byte, error) {

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)

}
