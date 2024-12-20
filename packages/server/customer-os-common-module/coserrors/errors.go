package coserrors

import (
	"github.com/pkg/errors"
)

var (
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
	ErrLinkedInUsed = errors.New("linkedin url is already used")
)

func SkipTracing(err error) bool {
	if err != nil {
		return false
	}

	// List of errors to be skipped from tracing
	errs := []error{
		ErrLinkedInUsed,
	}

	for _, e := range errs {
		if errors.Is(err, e) {
			return true
		}
	}

	return false
}
