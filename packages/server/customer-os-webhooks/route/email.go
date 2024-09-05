package route

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func AddEmailRoutes(ctx context.Context, route *gin.Engine, cfg *config.Config, services *service.Services) {
	route.POST("/sync/enrow/email",
		tracing.TracingEnhancer(ctx, "/sync/enrow/email"),
		syncEnrowEmailResponse(cfg, services))
}

func syncEnrowEmailResponse(cfg *config.Config, services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "syncEnrowEmailResponse", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		apiKeyHeader := c.Query("apiKey")
		if apiKeyHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing api key"})
			return
		}

		if apiKeyHeader != cfg.EnrowCallbackApiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		// Limit the size of the request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error reading request body"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var enrowResponseBody postgresentity.EnrowResponseBody
		if err = json.Unmarshal(body, &enrowResponseBody); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error unmarshalling request body"))
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		ctx, cancel := context.WithTimeout(ctx, common.Min1Duration)
		defer cancel()

		err = services.CommonServices.PostgresRepositories.CacheEmailEnrowRepository.AddResponse(ctx, enrowResponseBody.Id, enrowResponseBody.Qualification, string(body))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error adding enrow response"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing enrow response"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	}
}
