package pkg

import "github.com/gofiber/fiber/v2"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func ErrorResponse(c *fiber.Ctx, err *AppError) error {
	return c.Status(err.Code).JSON(fiber.Map{
		"error": err.Message,
	})
}
