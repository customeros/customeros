package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	_ "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
)

func InitMailboxesRoutes(r *gin.Engine, services *service.Services) {

	r.GET("/mailboxes",
		security.TenantUserContextEnhancer(security.TENANT, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API, security.WithCache(services.CommonServices.Cache)),
		getMailboxesHandler(services))

}

// @Accept  json
// @Produce  json
// @Success 200 {array} entity.TenantSettingsMailbox
// @Failure 401
// @Failure 500
// @Router /mailboxes [get]
// @Param   X-CUSTOMER-OS-API-KEY  header  string  true  "Authorization token"
func getMailboxesHandler(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /sequences", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		span.SetTag(tracing.SpanTagTenant, tenant)

		mailboxes, err := services.CommonServices.PostgresRepositories.TenantSettingsMailboxRepository.GetAll(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, mailboxes)
	}
}
