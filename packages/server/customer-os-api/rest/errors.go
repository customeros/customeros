package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
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
	ErrConflict = &ErrorResponse{
		BaseResponse: BuildBaseResponse(StatusError),
		Message:      "Entity already exists",
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

func SendError(c *gin.Context, span opentracing.Span, httpStatusCode int, err *ErrorResponse) {
	tracing.LogObjectAsJson(span, c.Request.URL.Path, err)
	c.JSON(httpStatusCode, err)
}
