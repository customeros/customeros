package errors

import "github.com/pkg/errors"

var (
	ErrAccessDenied      = errors.New("Access denied")
	ErrInvalidEntityType = errors.New("Invalid entity type")
	ErrMissingInput      = errors.New("Missing input")
	ErrNotFound          = errors.New("Not found")
)
