package commonerrors

import "github.com/pkg/errors"

var (
	ErrOperationNotAllowed = errors.New("Operation not allowed")
)
