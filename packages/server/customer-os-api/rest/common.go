package rest

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

// BaseResponse represents the standard API response structure
// @Description Standard response structure for API operations
type BaseResponse struct {
	// Status indicates the result of the operation ("success" or "error")
	Status string `json:"status" example:"success"`

	// Message provides additional information about the operation
	Message string `json:"message,omitempty" example:"Operation completed successfully"`
}

type HTTPContext struct {
	GinContext     *gin.Context
	ServiceContext *context.Context
	Span           opentracing.Span
	Services       *service.Services
	Tenant         string
}

func ValidateTenant(c *gin.Context, ctx context.Context, span opentracing.Span) string {
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)

	if tenant == "" {
		SendError(c, http.StatusUnauthorized, "API key invalid or expired")
		return ""
	}
	return tenant
}

func ValidateAndOpenCsvFile(c *gin.Context) (multipart.File, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		SendError(c, http.StatusBadRequest, "Failed to parse file")
		return nil, err
	}

	if header.Header.Get("Content-Type") != "text/csv" && !strings.HasSuffix(header.Filename, ".csv") {
		SendError(c, http.StatusBadRequest, "Invalid file type")
		return nil, errors.New("invalid file type")
	}

	return file, nil
}

func SendError(c *gin.Context, status int, message string) {
	c.JSON(status, BaseResponse{
		Status:  "error",
		Message: message,
	})
}

func SendSuccess(c *gin.Context, status int, message string) {
	c.JSON(status, BaseResponse{
		Status:  "success",
		Message: message,
	})
}

func CreateError(message string) BaseResponse {
	return BaseResponse{
		Status:  "error",
		Message: message,
	}
}

func CreateSuccess(message string) BaseResponse {
	return BaseResponse{
		Status:  "success",
		Message: message,
	}
}
