package route

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"io"
	"net/http"
	"time"
)

func AddIssueRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger) {
	route.POST("/sync/issues",
		cosHandler.TracingEnhancer(ctx, "/sync/issues"),
		commonService.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonService.CUSTOMER_OS_WEBHOOKS),
		syncIssuesHandler(services, log))
	route.POST("/sync/issue",
		cosHandler.TracingEnhancer(ctx, "/sync/issue"),
		commonService.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonService.CUSTOMER_OS_WEBHOOKS),
		syncIssueHandler(services, log))
}

func syncIssuesHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncIssues", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeCommon)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncIssues) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var issues []model.IssueData
		if err = json.Unmarshal(requestBody, &issues); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncIssues) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		if len(issues) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing issues in request"})
			return
		}

		// Context timeout, allocate per issue
		timeout := time.Duration(len(issues)) * utils.LongDuration
		if timeout > constants.RequestMaxTimeout {
			timeout = constants.RequestMaxTimeout
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err = services.IssueService.SyncIssues(ctx, issues)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncIssues) error in sync issues: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing issues"})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Message received successfully"})
		}
	}
}

func syncIssueHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncIssue", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeCommon)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncIssue) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var issue model.IssueData
		if err = json.Unmarshal(requestBody, &issue); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncIssue) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		// Context timeout, allocate per issue
		timeout := utils.LongDuration
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err = services.IssueService.SyncIssues(ctx, []model.IssueData{issue})
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncIssue) error in sync issue: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing issue"})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Message received successfully"})
		}
	}
}
