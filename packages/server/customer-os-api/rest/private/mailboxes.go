package private

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"

	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
)

func GetMailboxes(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "GET /mailboxes", c.Request.Header)
		defer span.Finish()

		tenant := c.Keys["TenantName"].(string)

		ctx = common.WithCustomContext(ctx, &common.CustomContext{
			Tenant: tenant,
		})

		span.SetTag(tracing.SpanTagTenant, tenant)

		mailboxes, err := s.Repositories.PostgresRepositories.TenantSettingsMailboxRepository.GetAll(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			c.Status(500)
			return
		}

		c.JSON(200, mailboxes)
	}
}
