package ws

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// UpgradeHandler returns a Fiber middleware that upgrades HTTP to WebSocket.
// Must be applied before the WebSocket handler.
func UpgradeHandler() fiber.Handler {

	return func(c *fiber.Ctx) error {

		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired

	}

}

// ConnectHandler returns the WebSocket connection handler.
// Expects userID to be set via query param: /ws?user_id=123
// In production, this should verify a token instead.
func ConnectHandler(hub *Hub) fiber.Handler {

	return websocket.New(func(c *websocket.Conn) {

		// Get user ID from query param (simplified — use token in production)
		userIDStr := c.Query("user_id", "0")
		var userID uint
		if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil || userID == 0 {
			log.Println("WS: missing user_id")
			c.Close()
			return
		}

		hub.Register(userID, c)
		defer hub.Unregister(userID, c)

		// Keep connection alive — read messages (ping/pong)
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}

	})

}
