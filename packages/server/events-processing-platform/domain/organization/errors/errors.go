package errors

import "github.com/pkg/errors"

var (
	ErrMissingPhoneNumberId = errors.New("missing phone number id")
	ErrMissingEmailId       = errors.New("missing email id")
	ErrMissingDomain        = errors.New("missing domain")
	ErrPhoneNumberNotFound  = errors.New("phone number not found")
	ErrEmailNotFound        = errors.New("email not found")
)
