package errors

import "github.com/pkg/errors"

var (
	ErrMissingEmailId      = errors.New("missing email id")
	ErrPhoneNumberNotFound = errors.New("phone number not found")
	ErrEmailNotFound       = errors.New("email not found")
)
