package validator

import "github.com/go-playground/validator/v10"

// Validate ...
var Validate *validator.Validate

func init() {
	Validate = validator.New()
}