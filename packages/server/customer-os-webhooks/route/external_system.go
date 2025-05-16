package route

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

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

func AddExternalSystemRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger, cache *commoncaches.Cache) {
	route.POST("/sync/external-system",
		RestTracingEnhancer(ctx, "/sync/external-system"),
		security.ApiKeyCheckerHTTP(services.PostgresRepository.TenantWebhookApiKeyRepository, services.Cfg.App.AppKey, security.WithCache(cache)),
		syncExternalSystemHandler(services, log))
}

func syncExternalSystemHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncExternalSystem", c.Request.Header)
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
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeMessages)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			spans.TraceError(err)
			log.Errorf("(SyncExternalSystem) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var externalSystemData model.ExternalSystemData
		if err = json.Unmarshal(requestBody, &externalSystemData); err != nil {
			spans.TraceError(err)
			log.Errorf("(SyncExternalSystem) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		ctx, cancel := context.WithTimeout(ctx, constants.Duration1Min)
		defer cancel()

		syncResult, err := services.ExternalSystemService.SyncExternalSystem(ctx, externalSystemData)
		if err != nil {
			spans.TraceError(err)
			log.Errorf("(SyncExternalSystem) error in sync external system: %s", err.Error())
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
