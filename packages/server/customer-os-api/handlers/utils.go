package handlers

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

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
)

// BaseResponse represents the standard API response structure
// @Description Standard response structure for API operations

func ValidateTenant(c *gin.Context, ctx context.Context, span opentracing.Span) string {
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)

	if tenant == "" {
		SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
		return ""
	}
	return tenant
}

func ValidateAndOpenCsvFile(c *gin.Context, span opentracing.Span) (multipart.File, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Unable to parse file"))
		return nil, err
	}

	if header.Header.Get("Content-Type") != "text/csv" && !strings.HasSuffix(header.Filename, ".csv") {
		SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Invalid file type"))
		return nil, errors.New("invalid file type")
	}

	return file, nil
}

func SendError(c *gin.Context, span opentracing.Span, httpStatusCode int, err *enum.ErrorResponse) {
	tracing.LogObjectAsJson(span, c.Request.URL.Path, err)
	c.JSON(httpStatusCode, err)
}
