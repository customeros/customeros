package leads_errors

import "errors"

var (
	ErrTenantMissing      = errors.New("Tenant not set on context")
	ErrUserIdMissing      = errors.New("UserID not set on context")
	ErrWebtrackerNotFound = errors.New("Webtracker not found")
	ErrWebtrackerExists   = errors.New("Webtracker already exists")
)
