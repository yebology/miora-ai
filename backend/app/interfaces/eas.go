// Package interfaces defines contracts for the EAS (Ethereum Attestation Service) client.
package interfaces

// IEASClient defines the contract for interacting with EAS on Base Sepolia.
// Used by the wallet service to publish trading reputation attestations on-chain.
type IEASClient interface {
	// Attest publishes a trading reputation attestation on-chain.
	// Returns the attestation UID (bytes32 hex) and transaction hash.
	Attest(recipient string, score uint8, recommendation string, totalTxns uint32, chain string) (uid string, txHash string, err error)

	// GetAttestation retrieves an attestation by its UID.
	// Returns the decoded attestation data.
	GetAttestation(uid string) (*AttestationData, error)
}

// AttestationData holds the decoded data from an on-chain EAS attestation.
type AttestationData struct {
	UID               string `json:"uid"`
	Schema            string `json:"schema"`
	Attester          string `json:"attester"`
	Recipient         string `json:"recipient"`
	Time              uint64 `json:"time"`
	ExpirationTime    uint64 `json:"expiration_time"`
	RevocationTime    uint64 `json:"revocation_time"`
	Revocable         bool   `json:"revocable"`
	Score             uint8  `json:"score"`
	Recommendation    string `json:"recommendation"`
	TotalTransactions uint32 `json:"total_transactions"`
	Chain             string `json:"chain"`
}
