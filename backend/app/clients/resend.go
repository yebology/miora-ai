// Resend client for sending email notifications.
// Uses Resend API (free tier: 100 emails/day).
//
// API docs: https://resend.com/docs/api-reference/emails/send-email
package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Resend sends emails via the Resend API.
type Resend struct {
	apiKey string
	from   string // Sender email (e.g. "Miora AI <alerts@miora.ai>")
}

// NewResend creates a new Resend client.
func NewResend(apiKey, fromEmail string) *Resend {

	return &Resend{
		apiKey: apiKey,
		from:   fromEmail,
	}

}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

type resendResponse struct {
	ID string `json:"id"`
}

// SendTradeAlert sends a trade notification email to the user.
func (r *Resend) SendTradeAlert(to, walletAddress, chain, tokenSymbol, direction, value string, liquidity, marketCap float64) error {

	action := "bought"
	actionColor := "#22c55e"
	if direction == "out" {
		action = "sold"
		actionColor = "#ef4444"
	}

	shortAddr := walletAddress
	if len(walletAddress) > 12 {
		shortAddr = walletAddress[:6] + "..." + walletAddress[len(walletAddress)-4:]
	}

	subject := fmt.Sprintf("🔔 %s %s %s on %s", shortAddr, action, tokenSymbol, chain)

	html := fmt.Sprintf(`
		<div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; max-width: 480px; margin: 0 auto; padding: 24px;">
			<h2 style="margin: 0 0 16px; font-size: 18px;">🔔 Trade Alert</h2>
			<div style="background: #f8f9fa; border-radius: 12px; padding: 20px; margin-bottom: 16px;">
				<p style="margin: 0 0 8px; color: #666; font-size: 13px;">Watched wallet</p>
				<p style="margin: 0 0 16px; font-family: monospace; font-size: 14px;">%s</p>
				<p style="margin: 0; font-size: 16px;">
					<span style="color: %s; font-weight: 600;">%s</span>
					<strong>%s %s</strong> on <strong>%s</strong>
				</p>
			</div>
			<table style="width: 100%%; font-size: 13px; color: #666;">
				<tr><td style="padding: 4px 0;">Amount</td><td style="text-align: right;">%s</td></tr>
				<tr><td style="padding: 4px 0;">Liquidity</td><td style="text-align: right;">$%s</td></tr>
				<tr><td style="padding: 4px 0;">Market Cap</td><td style="text-align: right;">$%s</td></tr>
			</table>
			<hr style="border: none; border-top: 1px solid #eee; margin: 16px 0;">
			<p style="margin: 0; font-size: 11px; color: #999;">
				Miora AI — Not financial advice. Trade at your own risk.
			</p>
		</div>`,
		walletAddress,
		actionColor, action, value, tokenSymbol, chain,
		value,
		formatEmailNumber(liquidity),
		formatEmailNumber(marketCap),
	)

	payload := resendRequest{
		From:    r.from,
		To:      []string{to},
		Subject: subject,
		HTML:    html,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal email: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend error %d: %s", resp.StatusCode, string(respBody))
	}

	var result resendResponse
	json.NewDecoder(resp.Body).Decode(&result)
	log.Printf("Email sent: %s → %s (id: %s)", subject, to, result.ID)

	return nil

}

// formatEmailNumber formats a number for email display.
func formatEmailNumber(n float64) string {
	if n >= 1_000_000_000 {
		return fmt.Sprintf("%.1fB", n/1_000_000_000)
	}
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", n/1_000_000)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%.0fk", n/1_000)
	}
	return fmt.Sprintf("%.0f", n)
}
