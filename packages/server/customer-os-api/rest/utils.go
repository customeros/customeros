package rest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

// BaseResponse represents the standard API response structure
// @Description Standard response structure for API operations
type BaseResponse struct {
	RequestID string `json:"requestId" example:"1234567890abcdef"`
	// Status indicates the result of the operation ("success" or "error")
	Status string `json:"status" example:"success"`
}

type HTTPContext struct {
	GinContext     *gin.Context
	ServiceContext *context.Context
	Span           opentracing.Span
	Services       *service.Services
	Tenant         string
}

type Status string

const (
	StatusError          Status = "error"
	StatusPartialSuccess Status = "partial success"
	StatusProcessing     Status = "processing"
	StatusSuccess        Status = "success"
	StatusWarning        Status = "warning"
)

const requestIDLength = 16

func BuildBaseResponse(status Status) BaseResponse {
	return BaseResponse{
		RequestID: generateRequestID(),
		Status:    string(status),
	}
}

func ValidateTenant(c *gin.Context, ctx context.Context, span opentracing.Span) string {
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)

	if tenant == "" {
		SendError(c, span, http.StatusUnauthorized, ErrInvalidAPIKey)
		return ""
	}
	return tenant
}

func ValidateAndOpenCsvFile(c *gin.Context, span opentracing.Span) (multipart.File, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		SendError(c, span, http.StatusBadRequest, ErrBadRequest.WithMessage("Unable to parse file"))
		return nil, err
	}

	if header.Header.Get("Content-Type") != "text/csv" && !strings.HasSuffix(header.Filename, ".csv") {
		SendError(c, span, http.StatusBadRequest, ErrBadRequest.WithMessage("Invalid file type"))
		return nil, errors.New("invalid file type")
	}

	return file, nil
}

func generateRequestID() string {
	bytes := make([]byte, requestIDLength/2)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp if crypto/rand fails
		return fmt.Sprintf("%x", time.Now().UnixNano())[:requestIDLength]
	}
	return hex.EncodeToString(bytes)
}
