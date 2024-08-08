package route

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
)

func AddLogEntryRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger, cache *commoncaches.Cache) {
	route.POST("/sync/log-entries",
		tracing.TracingEnhancer(ctx, "/sync/log-entries"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_WEBHOOKS, security.WithCache(cache)),
		syncLogEntriesHandler(services, log))
	route.POST("/sync/log-entry",
		tracing.TracingEnhancer(ctx, "/sync/log-entry"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_WEBHOOKS, security.WithCache(cache)),
		syncLogEntryHandler(services, log))
}

func syncLogEntriesHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncLogEntries", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeMessages)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncLogEntries) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var logEntries []model.LogEntryData
		if err = json.Unmarshal(requestBody, &logEntries); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncLogEntries) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		if len(logEntries) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing logEntries in request"})
			return
		}

		// Context timeout, allocate per logEntry
		timeout := time.Duration(len(logEntries)) * common.Min1Duration
		if timeout > constants.RequestMaxTimeout {
			timeout = constants.RequestMaxTimeout
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		syncResult, err := services.LogEntryService.SyncLogEntries(ctx, logEntries)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncLogEntries) error in sync logEntries: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing log entries"})
			}
		} else {
			c.JSON(http.StatusOK, syncResult)
		}
	}
}

func syncLogEntryHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncLogEntry", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeMessages)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncLogEntry) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var logEntry model.LogEntryData
		if err = json.Unmarshal(requestBody, &logEntry); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncLogEntries) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		// Context timeout, allocate per logEntry
		ctx, cancel := context.WithTimeout(ctx, common.Min1Duration)
		defer cancel()

		syncResult, err := services.LogEntryService.SyncLogEntries(ctx, []model.LogEntryData{logEntry})
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncLogEntry) error in sync logEntry: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing log entry"})
			}
		} else {
			c.JSON(http.StatusOK, syncResult)
		}
	}
}
