package utils

import "github.com/go-playground/validator/v10"

// GetValidator returns a new instance of the go-playground validator for input validation.
func GetValidator() *validator.Validate {

	return validator.New()

}
