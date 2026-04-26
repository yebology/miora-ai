package requests

// RegenerateInsight is the request body for POST /wallets/regenerate-insight.
type RegenerateInsight struct {
	Address      string `json:"address" validate:"required"`
	Chain        string `json:"chain" validate:"required"`
	Tone         string `json:"tone" validate:"required,oneof=simple eli5 custom"`
	CustomPrompt string `json:"custom_prompt"` // Required when tone is "custom", max 200 chars
}
