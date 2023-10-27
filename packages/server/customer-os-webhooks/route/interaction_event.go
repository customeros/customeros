package route

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"io"
	"net/http"
	"time"
)

func AddInteractionEventRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger) {
	route.POST("/sync/interaction-events",
		handler.TracingEnhancer(ctx, "/sync/interaction-events"),
		commonservice.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonservice.CUSTOMER_OS_WEBHOOKS),
		syncInteractionEventsHandler(services, log))
	route.POST("/sync/interaction-event",
		handler.TracingEnhancer(ctx, "/sync/interaction-event"),
		commonservice.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonservice.CUSTOMER_OS_WEBHOOKS),
		syncInteractionEventHandler(services, log))
}

func syncInteractionEventsHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncInteractionEvents", c.Request.Header)
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
			log.Errorf("(SyncInteractionEvents) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var interactionEvents []model.InteractionEventData
		if err = json.Unmarshal(requestBody, &interactionEvents); err != nil {
			log.Errorf("(SyncInteractionEvents) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		if len(interactionEvents) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing interactionEvents in request"})
			return
		}

		// Context timeout, allocate per interactionEvent
		timeout := time.Duration(len(interactionEvents)) * utils.LongDuration
		if timeout > constants.RequestMaxTimeout {
			timeout = constants.RequestMaxTimeout
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err = services.InteractionEventService.SyncInteractionEvents(ctx, interactionEvents)
		if err != nil {
			log.Errorf("(SyncInteractionEvents) error in sync interactionEvents: %s", err.Error())
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

func syncInteractionEventHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncInteractionEvent", c.Request.Header)
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
			log.Errorf("(SyncInteractionEvent) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var interactionEvent model.InteractionEventData
		if err = json.Unmarshal(requestBody, &interactionEvent); err != nil {
			log.Errorf("(SyncInteractionEvents) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		// Context timeout, allocate per interactionEvent
		timeout := utils.LongDuration
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err = services.InteractionEventService.SyncInteractionEvents(ctx, []model.InteractionEventData{interactionEvent})
		if err != nil {
			log.Errorf("(SyncInteractionEvent) error in sync interactionEvent: %s", err.Error())
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
