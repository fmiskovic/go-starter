package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validator responsibility is to validate incoming request objects.
type Validator struct {
	validate *validator.Validate
}

// NewValidator instantiate new Validator.
func NewValidator() Validator {
	return Validator{validate: validator.New()}
}

// Validate incoming request object and returns array of errors or empty.
func (v Validator) Validate(data interface{}) []string {
	errors := make([]string, 0)

	errs := v.validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'",
				err.Field(),
				err.Value(),
				err.Tag(),
			))
		}
	}

	return errors
}
