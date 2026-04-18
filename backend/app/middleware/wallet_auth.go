// Package middleware provides the wallet-based authentication middleware.
//
// Instead of Firebase, users authenticate by connecting their wallet (MetaMask).
// The frontend sends the wallet address in the X-Wallet-Address header.
// The middleware finds or creates a user record based on this address.
package middleware

import (
	"strings"

	"miora-ai/app/output"
	"miora-ai/constants"

	"github.com/gofiber/fiber/v2"
)

// WalletAuth verifies the wallet address from the X-Wallet-Address header.
// On success, sets "wallet_address" in c.Locals for handlers to use.
//
// Header format: X-Wallet-Address: 0x1234...
func WalletAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		walletAddress := c.Get("X-Wallet-Address")
		if walletAddress == "" {
			return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
		}

		// Normalize to lowercase
		walletAddress = strings.ToLower(walletAddress)

		// Basic validation: must start with 0x and be 42 chars
		if len(walletAddress) != 42 || !strings.HasPrefix(walletAddress, "0x") {
			return output.GetError(c, fiber.StatusBadRequest, "Invalid wallet address.")
		}

		c.Locals("wallet_address", walletAddress)
		return c.Next()
	}
}
