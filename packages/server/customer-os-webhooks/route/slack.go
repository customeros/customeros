package route

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
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
)

func AddSlackRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger, cache *commoncaches.Cache) {
	route.POST("/slack/channels",
		handler.TracingEnhancer(ctx, "/slack/channels"),
		commonservice.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.TenantApiKeyRepository, services.CommonServices.CommonRepositories.AppKeyRepository, commonservice.CUSTOMER_OS_WEBHOOKS, commonservice.WithCache(cache)),
		syncSlackChannelsHandler(services, log))
}

func syncSlackChannelsHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncSlackChannels", c.Request.Header)
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
			log.Errorf("(SyncSlackChannels) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var slackChannels []model.SlackChannelData
		if err = json.Unmarshal(requestBody, &slackChannels); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncSlackChannels) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		if len(slackChannels) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing slack channels in request"})
			return
		}

		for _, slackChannel := range slackChannels {
			err := services.CommonServices.SlackChannelService.StoreSlackChannel(tenant, slackChannel.ChannelId, nil)
			if err != nil {
				tracing.TraceErr(span, err)
				log.Errorf("(SyncSlackChannels) error in sync users: %s", err.Error())
				if errors.IsBadRequest(err) {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing users"})
				}
			} else {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			}
		}

		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncSlackChannels) error in sync users: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing users"})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	}
}
