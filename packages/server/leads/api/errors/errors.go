package api_errors

import (
	"fmt"
	"strings"
)

type MultiErrors struct {
	Errors map[string][]ErrorInfo
}

type ErrorInfo struct {
	Message  string
	RawError error
}

func NewMultiErrors() *MultiErrors {
	return &MultiErrors{
		Errors: make(map[string][]ErrorInfo),
	}
}

func (e *MultiErrors) Add(key, message string, err error) {
	e.Errors[key] = append(e.Errors[key], ErrorInfo{
		Message:  message,
		RawError: err,
	})
}

func (e *MultiErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

func (e *MultiErrors) Error() string {
	var parts []string
	for field, errors := range e.Errors {
		for _, err := range errors {
			parts = append(parts, fmt.Sprintf("%s: %s", field, err.Message))
		}
	}
	return strings.Join(parts, " | ")
}
