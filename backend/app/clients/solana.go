package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"miora-ai/app/interfaces"
)

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

func (a *AlchemySolana) GetTransfers(address string) ([]interfaces.TransferData, error) {
	url := fmt.Sprintf("https://solana-mainnet.g.alchemy.com/v2/%s", a.apiKey)

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

	transfers := make([]interfaces.TransferData, 0, len(result.Result))
	for _, tx := range result.Result {
		transfers = append(transfers, interfaces.TransferData{
			Hash:      tx.Signature,
			Timestamp: tx.BlockTime,
		})
	}

	return transfers, nil
}
