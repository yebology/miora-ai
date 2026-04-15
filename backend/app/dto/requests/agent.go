// Package requests contains agent request DTOs.
package requests

// UpdateAgentConfig is the request body for PUT /agent/config.
type UpdateAgentConfig struct {
	Budget        float64  `json:"budget" validate:"omitempty,gte=0"`
	MaxPerTrade   float64  `json:"max_per_trade" validate:"omitempty,gte=0"`
	RiskTolerance string   `json:"risk_tolerance" validate:"omitempty,oneof=low medium high"`
	MinScore      int      `json:"min_score" validate:"omitempty,gte=0,lte=100"`
	Conditions    []string `json:"conditions"`
}
