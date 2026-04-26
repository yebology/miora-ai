// Package ws provides WebSocket hub for real-time notifications.
//
// Hub manages WebSocket connections per user. When a watched wallet trades,
// the monitor service sends a message to the hub, which broadcasts to
// all connected clients of users who follow that wallet.
package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

// Message represents a notification sent to a user via WebSocket.
type Message struct {
	Type    string      `json:"type"`    // "wallet_trade", "alert", etc.
	Payload interface{} `json:"payload"` // Notification data
}

// Hub manages WebSocket connections grouped by user ID.
type Hub struct {
	// connections maps userID → set of WebSocket connections
	connections map[uint]map[*websocket.Conn]bool
	mu          sync.RWMutex
}

// NewHub creates a new WebSocket hub.
func NewHub() *Hub {

	return &Hub{
		connections: make(map[uint]map[*websocket.Conn]bool),
	}

}

// Register adds a WebSocket connection for a user.
func (h *Hub) Register(userID uint, conn *websocket.Conn) {

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connections[userID] == nil {
		h.connections[userID] = make(map[*websocket.Conn]bool)
	}
	h.connections[userID][conn] = true

	log.Printf("WS: user %d connected (total: %d)", userID, len(h.connections[userID]))

}

// Unregister removes a WebSocket connection for a user.
func (h *Hub) Unregister(userID uint, conn *websocket.Conn) {

	h.mu.Lock()
	defer h.mu.Unlock()

	if conns, ok := h.connections[userID]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(h.connections, userID)
		}
	}

	log.Printf("WS: user %d disconnected", userID)

}

// SendToUser sends a message to all connections of a specific user.
func (h *Hub) SendToUser(userID uint, msg Message) {

	h.mu.RLock()
	conns := h.connections[userID]
	h.mu.RUnlock()

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("WS: marshal error: %v", err)
		return
	}

	for conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("WS: write error for user %d: %v", userID, err)
			conn.Close()
			h.Unregister(userID, conn)
		}
	}

}

// SendToUsers sends a message to multiple users.
func (h *Hub) SendToUsers(userIDs []uint, msg Message) {

	for _, uid := range userIDs {
		h.SendToUser(uid, msg)
	}

}
