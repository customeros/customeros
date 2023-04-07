package errors

import "github.com/pkg/errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrEmailMissingId     = errors.New("email id is missing")
)
