package validator

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()
}

func GetValidator() *validator.Validate {
	if validate == nil {
		InitValidator()
	}
	return validate
}
