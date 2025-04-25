package api_errors

import "github.com/vektah/gqlparser/v2/gqlerror"

const (
	CodeNotFound     = "NOT_FOUND"
	CodeForbidden    = "FORBIDDEN"
	CodeBadInput     = "BAD_USER_INPUT"
	CodeInternal     = "INTERNAL_ERROR"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeExists       = "ALREADY_EXISTS"
)

// NewError creates a standardized GraphQL error
func NewError(message string, code string, extensions map[string]interface{}) *gqlerror.Error {
	if extensions == nil {
		extensions = make(map[string]interface{})
	}
	extensions["code"] = code

	return &gqlerror.Error{
		Message:    message,
		Extensions: extensions,
	}
}
