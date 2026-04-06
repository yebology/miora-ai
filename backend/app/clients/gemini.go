// Gemini client for generating natural language wallet insights.
// Uses Google Gemini API (free tier: 15 RPM for gemini-1.5-flash).
//
// API docs: https://ai.google.dev/api/generate-content
package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Gemini generates text using Google's Gemini API.
type Gemini struct {
	apiKey string
	model  string
}

// NewGemini creates a new Gemini client.
func NewGemini(apiKey string) *Gemini {

	return &Gemini{
		apiKey: apiKey,
		model:  "gemini-2.0-flash",
	}

}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []geminiPart `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Generate sends a prompt to Gemini and returns the text response.
func (g *Gemini) Generate(prompt string) (string, error) {

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		g.model, g.apiKey,
	)

	payload := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("gemini request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("gemini error %d: %s", resp.StatusCode, string(respBody))
	}

	var result geminiResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini: empty response")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil

}
