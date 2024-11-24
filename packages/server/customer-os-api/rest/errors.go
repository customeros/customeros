package rest

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	BaseResponse
	Message string `json:"message,omitempty"`
}

var (
	ErrBadRequest = &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      "Invalid request format",
	}
	ErrInvalidAPIKey = &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      "API key is invaild or expired",
	}
	ErrNotFound = &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      "Resource not found",
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

func SendError(c *gin.Context, httpStatusCode int, err *ErrorResponse) {
	c.JSON(httpStatusCode, err)
}
