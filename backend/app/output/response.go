// Package output provides standardized API response functions.
//
// Every API response uses the ApiResponse envelope:
//
//	Success: { "status": "success", "message": "...", "data": { ... } }
//	Error:   { "status": "error",   "message": "..." }
package output

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// ApiResponse is the standard response envelope for all API endpoints.
type ApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// GetSuccess returns a 200 response with status "success".
// Pass nil for data if there's nothing to return.
func GetSuccess(c *fiber.Ctx, message string, data interface{}) error {

	return c.Status(fiber.StatusOK).JSON(ApiResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})

}

// GetError logs the error and returns a response with status "error".
func GetError(c *fiber.Ctx, code int, message string) error {

	log.Printf("[ERROR] %d — %s", code, message)

	return c.Status(code).JSON(ApiResponse{
		Status:  "error",
		Message: message,
	})

}
