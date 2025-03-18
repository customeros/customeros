package ai

import "fmt"

type AIError struct {
	Message string
	Retry   bool
	Cause   error
}

func (e *AIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AIError) Unwrap() error {
	return e.Cause
}

func (s *aiService) NewRetryableError(message string, err error) *AIError {
	return &AIError{
		Message: message,
		Retry:   true,
		Cause:   err,
	}
}

func (s *aiService) NewErrorNoRetry(message string, err error) *AIError {
	return &AIError{
		Message: message,
		Retry:   false,
		Cause:   err,
	}
}

func (s *aiService) IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	e, ok := err.(*AIError)
	if ok {
		return e.Retry
	}
	return false
}
