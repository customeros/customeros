package errors

import "github.com/pkg/errors"

var (
	ErrMissingDomain       = errors.New("missing domain")
	ErrPhoneNumberNotFound = errors.New("phone number not found")
	ErrEmailNotFound       = errors.New("email not found")
)
