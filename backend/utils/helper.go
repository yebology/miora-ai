package utils

import (
	"miora-ai/constants"
	"miora-ai/pkg"

	"github.com/gofiber/fiber/v2"
)

// ParseAndValidateBody parses the JSON request body into dst and validates it.
// Returns *pkg.AppError if parsing or validation fails, nil on success.
func ParseAndValidateBody(c *fiber.Ctx, dst interface{}) *pkg.AppError {

	if err := c.BodyParser(dst); err != nil {
		return pkg.ErrBadReq(constants.InvalidRequest)
	}

	if err := GetValidator().Struct(dst); err != nil {
		return pkg.ErrBadReq(constants.InvalidData)
	}

	return nil

}
