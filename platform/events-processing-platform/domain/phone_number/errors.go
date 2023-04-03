package phone_number

import "github.com/pkg/errors"

var (
	ErrMissingTenant            = errors.New("missing tenant")
	ErrPhoneNumberAlreadyExists = errors.New("phone number already created")
)
