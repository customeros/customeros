package phone_number

import "github.com/pkg/errors"

var (
	ErrPhoneNumberAlreadyExists = errors.New("phone number already created")
)
