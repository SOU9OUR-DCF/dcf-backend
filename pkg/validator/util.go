package validator

import (
	"github.com/go-playground/validator/v10"
)

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates a struct
func Validate(i interface{}) error {
	return validate.Struct(i)
}

// ValidateVar validates a variable
func ValidateVar(field interface{}, tag string) error {
	return validate.Var(field, tag)
}