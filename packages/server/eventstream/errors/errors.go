package inbox_errors

import "errors"

var (
	ErrTenantMissing = errors.New("Tenant not set on context")
	ErrUserIdMissing = errors.New("UserID not set on context")
)
