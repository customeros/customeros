package error

import "github.com/pkg/errors"

var (
	ErrLocationAlreadyExists = errors.New("location already exists")
)
