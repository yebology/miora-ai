// Package middleware provides the x402 payment verification middleware.
//
// x402 is a protocol that uses HTTP 402 (Payment Required) to gate API access
// behind USDC micropayments on Base. When a request lacks payment, the middleware
// returns 402 with payment requirements. When payment is provided via the
// X-PAYMENT header, it verifies the payment through a facilitator service.
//
// Reference: https://github.com/coinbase/x402
package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"miora-ai/constants"

	"github.com/gofiber/fiber/v2"
	x402 "github.com/mark3labs/x402-go"
)

// X402Config holds the configuration for x402 payment middleware.
type X402Config struct {
	RecipientAddress string // USDC receiving address on Base Sepolia
	PriceUSDC        string // Price per request in USDC (e.g. "0.01")
	FacilitatorURL   string // x402 facilitator URL for payment verification
}

// X402Middleware returns a Fiber middleware that gates access behind x402 USDC payments.
// If no payment header is present, returns 402 with payment requirements.
// If payment is present, verifies it through the facilitator before allowing access.
func X402Middleware(cfg X402Config) fiber.Handler {
	if cfg.FacilitatorURL == "" {
		cfg.FacilitatorURL = constants.X402FacilitatorURL
	}

	// Build payment requirement using x402-go helper
	requirement, err := x402.NewUSDCPaymentRequirement(x402.USDCRequirementConfig{
		Chain:            x402.BaseSepolia,
		Amount:           cfg.PriceUSDC,
		RecipientAddress: cfg.RecipientAddress,
	})
	if err != nil {
		// If config is invalid, return a middleware that always rejects
		return func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("x402 configuration error: %v", err),
			})
		}
	}

	requirementsJSON, _ := json.Marshal([]x402.PaymentRequirement{requirement})

	return func(c *fiber.Ctx) error {
		// Check for X-PAYMENT header
		paymentHeader := c.Get("X-PAYMENT")
		if paymentHeader == "" {
			// No payment — return 402 with requirements
			c.Set("Content-Type", "application/json")
			return c.Status(constants.HTTPStatusPaymentRequired).JSON(fiber.Map{
				"status":  "error",
				"message": "Payment required. Include X-PAYMENT header with x402 payment payload.",
				"x402": fiber.Map{
					"version":      "1",
					"requirements": json.RawMessage(requirementsJSON),
				},
			})
		}

		// Verify payment through facilitator
		if err := verifyPayment(cfg.FacilitatorURL, paymentHeader, requirementsJSON); err != nil {
			return c.Status(constants.HTTPStatusPaymentRequired).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("Payment verification failed: %v", err),
			})
		}

		// Payment verified — proceed to handler
		return c.Next()
	}
}

// verifyPayment sends the payment payload to the x402 facilitator for verification.
func verifyPayment(facilitatorURL, paymentHeader string, requirements []byte) error {
	payload, _ := json.Marshal(map[string]any{
		"payment":      paymentHeader,
		"requirements": json.RawMessage(requirements),
	})

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(facilitatorURL+"/verify", "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("facilitator request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("facilitator rejected payment: %s", string(body))
	}

	// Parse response to check if valid
	var result struct {
		Valid bool `json:"valid"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("invalid facilitator response: %w", err)
	}

	if !result.Valid {
		return fmt.Errorf("payment is not valid")
	}

	return nil
}
