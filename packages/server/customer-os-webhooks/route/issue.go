package route

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	commoncaches "github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/customer-os-webhooks/constants"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/errors"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/model"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/service"
)

func AddIssueRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger, cache *commoncaches.Cache) {
	route.POST("/sync/issues",
		RestTracingEnhancer(ctx, "/sync/issues"),
		security.ApiKeyCheckerHTTP(services.PostgresRepository.TenantWebhookApiKeyRepository, services.Cfg.App.AppKey, security.WithCache(cache)),
		syncIssuesHandler(services, log))
	route.POST("/sync/issue",
		RestTracingEnhancer(ctx, "/sync/issue"),
		security.ApiKeyCheckerHTTP(services.PostgresRepository.TenantWebhookApiKeyRepository, services.Cfg.App.AppKey, security.WithCache(cache)),
		syncIssueHandler(services, log))
}

func syncIssuesHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncIssues", c.Request.Header)
		defer spans.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant:    tenant,
			AppSource: constants.AppSourceCustomerOsWebhooks,
		})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeCommon)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			spans.TraceError(err)
			log.Errorf("(SyncIssues) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var issues []model.IssueData
		if err = json.Unmarshal(requestBody, &issues); err != nil {
			spans.TraceError(err)
			log.Errorf("(SyncIssues) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		if len(issues) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing issues in request"})
			return
		}

		// Context timeout, allocate per issue
		timeout := time.Duration(len(issues)) * constants.Duration1Min
		if timeout > constants.RequestMaxTimeout {
			timeout = constants.RequestMaxTimeout
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		syncResult, err := services.IssueService.SyncIssues(ctx, issues)
		if err != nil {
			spans.TraceError(err)
			log.Errorf("(SyncIssues) error in sync issues: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing issues"})
			}
		} else {
			c.JSON(http.StatusOK, syncResult)
		}
	}
}

func syncIssueHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncIssue", c.Request.Header)
		defer spans.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		span.SetTag(tracing.SpanTagTenant, tenant)
		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeCommon)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			spans.TraceError(err)
			log.Errorf("(SyncIssue) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var issue model.IssueData
		if err = json.Unmarshal(requestBody, &issue); err != nil {
			spans.TraceError(err)
			log.Errorf("(SyncIssue) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		// Context timeout, allocate per issue
		ctx, cancel := context.WithTimeout(ctx, constants.Duration1Min)
		defer cancel()

		syncResult, err := services.IssueService.SyncIssues(ctx, []model.IssueData{issue})
		if err != nil {
			spans.TraceError(err)
			log.Errorf("(SyncIssue) error in sync issue: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing issue"})
			}
		} else {
			c.JSON(http.StatusOK, syncResult)
		}
	}
}
