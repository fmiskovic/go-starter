package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() Validator {
	return Validator{validate: validator.New()}
}

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
