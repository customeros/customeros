package errors

import "github.com/pkg/errors"

var (
	ErrInvalidEntityType = errors.New("Invalid entity type")
	ErrMissingInput      = errors.New("Missing input")
)
