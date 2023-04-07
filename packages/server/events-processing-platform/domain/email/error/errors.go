package error

import "github.com/pkg/errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)
