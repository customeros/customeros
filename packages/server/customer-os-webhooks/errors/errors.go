package errors

import "github.com/pkg/errors"

var (
	ErrTenantNotValid            = errors.New("tenant not valid")
	ErrMissingExternalSystem     = errors.New("missing external system")
	ErrExternalSystemNotAccepted = errors.New("external system not accepted")
)
