package webhooks

import (
	"net/http"
	"strings"

	commonEnum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (h *WebhookHandler) HandleWebhook(flowsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "Webhooks")
		defer spans.Finish()

		spans.LogKV("param.tenantHash", c.Param("tenantHash"))

		tenant, err := h.services.Repositories.PostgresRepositories.TenantRepository.GetTenantByHashId(ctx, c.Param("tenantHash"))
		if err != nil {
			err := errors.Wrap(err, "Unable to identify tenant")
			spans.TraceError(err)
			h.responseHandler.HandleError(c, http.StatusUnauthorized, nil)
			return
		}

		webhookPath := strings.TrimPrefix(c.Request.URL.Path, flowsPath)
		integration, err := h.services.CommonServices.WebhookService.GetIntegrationFromWebhookPath(ctx, tenant, webhookPath)
		if err != nil {
			message := "Webhook not found"
			h.responseHandler.HandleError(c, http.StatusNotFound, &message)
			return
		}
		spans.LogKV("result.integration", integration.String())

		switch integration {
		case commonEnum.SourceCalCom:
			h.integrationsHandler.CalDotCom(c, tenant)
		// todo
		case commonEnum.SourceFathom:
			h.integrationsHandler.FathomZapier(c, tenant)
		case commonEnum.SourceGrain:
			h.integrationsHandler.GrainZapier(c, tenant)
		case commonEnum.SourcePostmark:
			h.integrationsHandler.PostmarkInboundEmail(c)
		// todo
		default:
			err := errors.New("Unsupported integration")
			spans.TraceError(err)
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}
	}
}
