// Package responses contains the reputation API response DTOs.
package responses

// Reputation represents the on-chain trading reputation data for a wallet.
type Reputation struct {
	Address           string `json:"address"`
	Chain             string `json:"chain"`
	Score             uint8  `json:"score"`
	Recommendation    string `json:"recommendation"`
	TotalTransactions uint32 `json:"total_transactions"`
	AttestationUID    string `json:"attestation_uid"`
	AttestationTxHash string `json:"attestation_tx_hash,omitempty"`
	Attester          string `json:"attester"`
	Timestamp         uint64 `json:"timestamp"`
	ExplorerURL       string `json:"explorer_url"`
}
