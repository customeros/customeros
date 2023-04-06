package errors

import "github.com/pkg/errors"

var (
	ErrMissingPhoneNumberId = errors.New("missing phone number id")
	ErrPhoneNumberNotFound  = errors.New("phone number not found")
)
