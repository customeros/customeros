package errors

import "github.com/pkg/errors"

var (
	ErrTenantNotValid            = errors.New("wrong tenant")
	ErrMissingExternalSystem     = errors.New("missing external system")
	ErrExternalSystemNotAccepted = errors.New("external system not accepted")
)

func IsBadRequest(err error) bool {
	return errors.Is(err, ErrTenantNotValid) ||
		errors.Is(err, ErrMissingExternalSystem) ||
		errors.Is(err, ErrExternalSystemNotAccepted)
}
