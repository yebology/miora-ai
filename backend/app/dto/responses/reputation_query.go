// Package responses contains the x402-protected reputation query response DTO.
package responses

// ReputationQuery is the response for the x402-protected reputation query endpoint.
// This is a simplified version of Reputation — only the essential scoring data.
type ReputationQuery struct {
	Address        string `json:"address"`
	Score          uint8  `json:"score"`
	Recommendation string `json:"recommendation"`
	Chain          string `json:"chain"`
	AttestationUID string `json:"attestation_uid,omitempty"`
	ExplorerURL    string `json:"explorer_url,omitempty"`
}
