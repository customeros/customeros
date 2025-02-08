package coserrors

import (
	"github.com/pkg/errors"
)

var (
	// context errors
	ErrUserEmailNotSet = errors.New("User email not set")
	ErrTenantNotSet    = errors.New("Tenant not set")

	ErrAccessDenied        = errors.New("Access denied")
	ErrInvalidEntityType   = errors.New("Invalid entity type")
	ErrNotSupported        = errors.New("Not supported")
	ErrConnectionTimeout   = errors.New("Connection timeout")
	ErrOperationNotAllowed = errors.New("Operation not allowed")

	// domain errors
	ErrDomainUnavailable         = errors.New("domain unavailable")
	ErrDomainPremium             = errors.New("domain is premium")
	ErrDomainPriceExceeded       = errors.New("domain price exceeds the maximum allowed price")
	ErrDomainPriceNotFound       = errors.New("domain price not found")
	ErrDomainConfigurationFailed = errors.New("domain configuration failed")
	ErrDomainNotFound            = errors.New("domain not found")

	// mailbox errors
	ErrMailboxExists = errors.New("mailbox already exists")

	// validation errors
	ErrLinkedInUsed       = errors.New("linkedin url is already used")
	ErrEmailUsed          = errors.New("Email is already used")
	ErrCannotIdentifyUser = errors.New("Cannot identify CustomerOS user")

	// Capability errors
	ErrCapabilityDomainMissing         = errors.New("Missing domain")
	ErrCapabilityContactMissing        = errors.New("Missing contact email")
	ErrCapabilityHostnameNotConfigured = errors.New("Hostname not configured")
)

func SkipTracing(err error) bool {
	if err == nil {
		return true
	}

	// List of errors to be skipped from tracing
	errs := []error{
		ErrLinkedInUsed,
		ErrEmailUsed,
		ErrCapabilityDomainMissing,
		ErrCapabilityContactMissing,
		ErrCapabilityHostnameNotConfigured,
	}

	for _, e := range errs {
		if errors.Is(err, e) {
			return true
		}
	}

	return false
}
