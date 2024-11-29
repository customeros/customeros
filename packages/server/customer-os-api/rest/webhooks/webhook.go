package webhooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func Webhooks(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Webhooks", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant, err := s.CommonServices.PostgresRepositories.TenantRepository.GetTenantFromHashID(ctx, c.Param("tenantId"))
		if err != nil {
			err := errors.Wrap(err, "Unable to identify tenant")
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
			return
		}

		integrationId := c.Param("integrationId")
		integration, err := s.WebhookService.LookupIntegration(integrationId)
		if err != nil {
			err := errors.Wrap(err, "Unable to lookup integration")
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusNotFound, rest.ErrNotFound)
			return
		}

		// build http context
		httpContext := rest.HTTPContext{
			GinContext:     c,
			ServiceContext: &ctx,
			Span:           span,
			Services:       s,
			Tenant:         tenant,
		}

		switch integration {
		case service.IntegrationCalCom:
		// todo
		case service.IntegrationFathom:
		// todo
		case service.IntegrationGrain:
		// todo
		case service.IntegrationPostmark:
			PostmarkInboundEmail(httpContext)
		// todo
		default:
			rest.SendError(c, span, http.StatusNotFound, rest.ErrNotFound)
			return
		}
	}
}
