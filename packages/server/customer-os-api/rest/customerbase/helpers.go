package customerbase

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
)

func validateTenant(c *gin.Context, ctx context.Context, span opentracing.Span) string {
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)

	if tenant == "" {
		sendError(c, http.StatusUnauthorized, "API key invalid or expired")
		return ""
	}
	return tenant
}

func validateAndOpenFile(c *gin.Context) (multipart.File, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		sendError(c, http.StatusBadRequest, "Failed to parse file")
		return nil, err
	}

	if header.Header.Get("Content-Type") != "text/csv" && !strings.HasSuffix(header.Filename, ".csv") {
		sendError(c, http.StatusBadRequest, "Invalid file type")
		return nil, errors.New("invalid file type")
	}

	return file, nil
}

func sendError(c *gin.Context, status int, message string) {
	c.JSON(status, BaseResponse{
		Status:  "error",
		Message: message,
	})
}

func createError(message string) BaseResponse {
	return BaseResponse{
		Status:  "error",
		Message: message,
	}
}

func createSuccess(message string) BaseResponse {
	return BaseResponse{
		Status:  "success",
		Message: message,
	}
}
