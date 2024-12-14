package enum

type ErrorResponse struct {
	BaseResponse
	Message string `json:"message,omitempty"`
}

var (
	ErrBadRequest = &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      "Invalid request format",
	}
	ErrConflict = &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      "Entity already exists",
	}
	ErrForbidden = &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      "Forbidden",
	}
	ErrInternalServer = &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      "Unable to process request",
	}
	ErrInvalidAPIKey = &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      "API key is invaild or expired",
	}
	ErrNotFound = &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      "Resource not found",
	}
	ErrUnauthorized = &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      "Tenant not authorized for this service",
	}
	ErrUnsupportedContentType = &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      "Unsupported Content-Type",
	}
	// Add more common errors as needed
)

func (e *ErrorResponse) WithMessage(message string) *ErrorResponse {
	return &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      message,
	}
}
