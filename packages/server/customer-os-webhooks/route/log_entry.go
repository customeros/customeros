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

func AddLogEntryRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger) {
	route.POST("/sync/log-entries",
		cosHandler.TracingEnhancer(ctx, "/sync/log-entries"),
		commonService.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonService.CUSTOMER_OS_WEBHOOKS),
		syncLogEntriesHandler(services, log))
	route.POST("/sync/log-entry",
		cosHandler.TracingEnhancer(ctx, "/sync/log-entry"),
		commonService.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonService.CUSTOMER_OS_WEBHOOKS),
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
			log.Errorf("(SyncLogEntries) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var logEntries []model.LogEntryData
		if err = json.Unmarshal(requestBody, &logEntries); err != nil {
			log.Errorf("(SyncLogEntries) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		if len(logEntries) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing logEntries in request"})
			return
		}

		// Context timeout, allocate per logEntry
		timeout := time.Duration(len(logEntries)) * utils.LongDuration
		if timeout > constants.RequestMaxTimeout {
			timeout = constants.RequestMaxTimeout
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err = services.LogEntryService.SyncLogEntries(ctx, logEntries)
		if err != nil {
			log.Errorf("(SyncLogEntries) error in sync logEntries: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing log entries"})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Message received successfully"})
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
			log.Errorf("(SyncLogEntry) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var logEntry model.LogEntryData
		if err = json.Unmarshal(requestBody, &logEntry); err != nil {
			log.Errorf("(SyncLogEntries) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		// Context timeout, allocate per logEntry
		timeout := utils.LongDuration
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err = services.LogEntryService.SyncLogEntries(ctx, []model.LogEntryData{logEntry})
		if err != nil {
			log.Errorf("(SyncLogEntry) error in sync logEntry: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing log entry"})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Message received successfully"})
		}
	}
}
