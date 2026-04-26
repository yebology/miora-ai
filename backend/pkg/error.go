// Package pkg provides shared utilities across the application.
package pkg

import "github.com/gofiber/fiber/v2"

// AppError represents a structured application error with HTTP status code.
// Services return *AppError instead of Go's error for consistent error handling.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {

	return e.Message

}

// ErrInternal returns a 500 Internal Server Error.
func ErrInternal() *AppError {

	return &AppError{Code: fiber.StatusInternalServerError, Message: "Internal server error."}

}

// ErrNotFound returns a 404 Not Found with a custom message.
func ErrNotFound(message string) *AppError {

	return &AppError{Code: fiber.StatusNotFound, Message: message}

}

// ErrBadReq returns a 400 Bad Request with a custom message.
func ErrBadReq(message string) *AppError {

	return &AppError{Code: fiber.StatusBadRequest, Message: message}

}

// ErrForbidden returns a 403 Forbidden with a custom message.
func ErrForbidden(message string) *AppError {

	return &AppError{Code: fiber.StatusForbidden, Message: message}

}

// ErrUnauthorized returns a 401 Unauthorized with a custom message.
func ErrUnauthorized(message string) *AppError {

	return &AppError{Code: fiber.StatusUnauthorized, Message: message}

}

// ErrConflict returns a 409 Conflict with a custom message.
func ErrConflict(message string) *AppError {

	return &AppError{Code: fiber.StatusConflict, Message: message}

}

// ErrUnexpected returns a custom status code with a custom message.
func ErrUnexpected(code int, message string) *AppError {

	return &AppError{Code: code, Message: message}

}
