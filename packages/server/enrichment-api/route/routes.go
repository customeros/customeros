package route

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/service"
	"net/http"
)

func RegisterRoutes(ctx context.Context, r *gin.Engine, services *service.Services) {
	r.GET("/tbd",
		handler.TracingEnhancer(ctx, "/tbd"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.VALIDATION_API),
		enrichPerson(services))
}

func enrichPerson(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "enrichPerson", c.Request.Header)
		defer span.Finish()

		// return 200 for now
		c.JSON(http.StatusOK, "OK")
	}
}
