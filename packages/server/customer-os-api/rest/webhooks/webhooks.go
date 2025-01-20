package webhooks

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	commonEnum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/rest/integrations"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type WebhookHandler struct {
	services            *cosapi_services.Services
	responseHandler     *response.Response
	integrationsHandler *integrations.IntegrationHandler
}

func NewWebhookHandler(
	services *cosapi_services.Services,
	responseHandler *response.Response,
	integrationsHandler *integrations.IntegrationHandler,
) *WebhookHandler {
	return &WebhookHandler{
		services:            services,
		responseHandler:     responseHandler,
		integrationsHandler: integrationsHandler,
	}
}

type CreateWebhookRequest struct {
	Integration string `json:"integration"`
}

type CreateWebhookRecord struct {
	URL         string `json:"url"`
	Integration string `json:"integration"`
	Secret      string `json:"secret"`
}

type CreateWebhookResponse struct {
	Hook CreateWebhookRecord `json:"hook"`
}

func (h *WebhookHandler) CreateWebhook(baseURL, flowsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateWebhook", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		var req CreateWebhookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			message := "Missing parameter: intergration"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		integration, err := h.services.WebhookService.GetIntegration(strings.ToLower(req.Integration))
		if err != nil {
			message := "Invalid integration value"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		webhookPath, secret, err := h.services.WebhookService.CreateIntegrationWebhook(ctx, tenant, integration)
		if err != nil {
			message := "Unable to create webhook"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		record := CreateWebhookRecord{
			URL:         fmt.Sprintf("%s%s/%s", baseURL, flowsPath, webhookPath),
			Integration: integration.String(),
			Secret:      secret,
		}

		h.responseHandler.HandleSuccess(c, CreateWebhookResponse{
			Hook: record,
		})
	}
}

type NoActiveWebhooks struct {
	Message string `json:"message"`
}

type OneActiveWebhook struct {
	Hook ActiveWebhookRecord `json:"hook"`
}

type ActiveWebhooksResponse struct {
	Hooks []ActiveWebhookRecord `json:"hooks"`
}

type ActiveWebhookRecord struct {
	URL         string    `json:"url"`
	Integration string    `json:"integration"`
	CreatedAt   time.Time `json:"createdAt"`
	Active      bool      `json:"active"`
}

func (h *WebhookHandler) GetActiveWebhooks(baseURL, flowsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetActiveWebhooks", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		webhooks, err := h.services.Repositories.PostgresRepositories.WebhooksRepository.FindAll(ctx)
		if err != nil {
			err = fmt.Errorf("Unable to lookup active webhooks for %s: %v", tenant, err)
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		if len(*webhooks) == 0 {
			h.responseHandler.HandleSuccess(c, NoActiveWebhooks{
				Message: "No active webhooks",
			})
			return
		}

		results := make([]ActiveWebhookRecord, len(*webhooks))

		for _, webhook := range *webhooks {
			record := ActiveWebhookRecord{
				URL:         fmt.Sprintf("%s%s/%s", baseURL, flowsPath, webhook.WebhookPath),
				Integration: webhook.Integration,
				CreatedAt:   webhook.CreatedAt,
				Active:      webhook.Enabled,
			}
			results = append(results, record)
		}

		if len(*webhooks) == 1 {
			h.responseHandler.HandleSuccess(c, OneActiveWebhook{
				Hook: results[0],
			})
			return
		}

		h.responseHandler.HandleSuccess(c, ActiveWebhooksResponse{
			Hooks: results,
		})
	}
}

func (h *WebhookHandler) RotateWebhook(baseURL, flowsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RotateWebhook", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}
		tenantId := c.Param("tenantId")
		validTenant, err := h.services.WebhookService.ValidateTenantId(ctx, tenant, tenantId)
		if err != nil {
			message := "Unable to verify webhook ownership"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}
		if !validTenant {
			h.responseHandler.HandleError(c, http.StatusUnauthorized, nil)
			return
		}

		// Lookup integration
		webhookPath := strings.TrimSuffix(c.Request.URL.Path, "/rotate")
		webhookPath = strings.TrimPrefix(webhookPath, flowsPath)
		integration, err := h.services.WebhookService.GetIntegrationFromWebhookPath(ctx, tenant, webhookPath)
		if err != nil || integration == commonEnum.SourceUnknown {
			message := "Unable to identify webhook"
			h.responseHandler.HandleError(c, http.StatusNotFound, &message)
			return
		}

		// Call create to rotate webhook as it will automatically handle rotation
		webhookPath, secret, err := h.services.WebhookService.CreateIntegrationWebhook(ctx, tenant, integration)
		if err != nil {
			message := "Unable to rotate webhook"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		record := CreateWebhookRecord{
			URL:         fmt.Sprintf("%s%s/%s", baseURL, flowsPath, webhookPath),
			Integration: integration.String(),
			Secret:      secret,
		}

		h.responseHandler.HandleSuccess(c, CreateWebhookResponse{
			Hook: record,
		})
	}
}

func (h *WebhookHandler) DeactivateWebhook() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RotateWebhook", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}
		tenantId := c.Param("tenantId")
		validTenant, err := h.services.WebhookService.ValidateTenantId(ctx, tenant, tenantId)
		if err != nil {
			message := "Unable to verify webhook ownership"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}
		if !validTenant {
			h.responseHandler.HandleError(c, http.StatusUnauthorized, nil)
			return
		}

		// Call to deactivate webhook
		deactErr := h.services.WebhookService.DeactivateWebhook(ctx, strings.TrimSuffix(c.Request.URL.Path, "/rotate"))
		if deactErr != nil {
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		h.responseHandler.HandleSuccess(c, NoActiveWebhooks{
			Message: "Webhook successfully deactivated",
		})
	}
}

func (h *WebhookHandler) HandleWebhook(flowsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Webhooks", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant, err := h.services.Repositories.PostgresRepositories.TenantRepository.GetTenant(ctx, c.Param("tenantId"))
		if err != nil {
			err := errors.Wrap(err, "Unable to identify tenant")
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusUnauthorized, nil)
			return
		}

		webhookPath := strings.TrimPrefix(c.Request.URL.Path, flowsPath)
		integration, err := h.services.WebhookService.GetIntegrationFromWebhookPath(ctx, tenant, webhookPath)
		if err != nil {
			message := "Webhook not found"
			h.responseHandler.HandleError(c, http.StatusNotFound, &message)
			return
		}

		switch integration {
		case commonEnum.SourceCalCom:
			h.integrationsHandler.CalDotCom(c)
		// todo
		case commonEnum.SourceFathom:
			h.integrationsHandler.FathomZapier(c)
		case commonEnum.SourceGrain:
			h.integrationsHandler.GrainZapier(c)
		case commonEnum.SourcePostmark:
			h.integrationsHandler.PostmarkInboundEmail(c)
		// todo
		default:
			err := errors.New("Unsupported integration")
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}
	}
}
