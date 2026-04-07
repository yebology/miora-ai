// Package middleware provides HTTP middleware for Fiber.
package middleware

import (
	"context"
	"strings"

	"miora-ai/app/output"
	"miora-ai/constants"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
)

// FirebaseAuth verifies the Firebase ID token from the Authorization header.
// On success, sets "firebase_uid", "firebase_email", "firebase_name", "firebase_avatar" in c.Locals.
//
// Header format: Authorization: Bearer <firebase_id_token>
func FirebaseAuth(credentialsFile string) fiber.Handler {

	// Initialize Firebase Auth client once
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		panic("failed to initialize firebase: " + err.Error())
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		panic("failed to initialize firebase auth: " + err.Error())
	}

	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
		}

		decoded, err := verifyToken(authClient, token)
		if err != nil {
			return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
		}

		// Set user info in locals for handlers
		c.Locals("firebase_uid", decoded.UID)
		c.Locals("firebase_email", getClaimString(decoded, "email"))
		c.Locals("firebase_name", getClaimString(decoded, "name"))
		c.Locals("firebase_avatar", getClaimString(decoded, "picture"))

		return c.Next()

	}

}

func verifyToken(client *auth.Client, token string) (*auth.Token, error) {

	return client.VerifyIDToken(context.Background(), token)

}

func getClaimString(token *auth.Token, key string) string {

	if val, ok := token.Claims[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""

}
