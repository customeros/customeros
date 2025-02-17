package webhooks

import (
	"fmt"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
	"strings"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	commonEnum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
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

func (h *WebhookHandler) CreateWebhook(baseURL, apiPath string) gin.HandlerFunc {
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
			message := "Missing parameter: integration"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		integration := h.services.CommonServices.WebhookService.GetIntegration(strings.ToLower(req.Integration))
		if integration == enum.SourceUnknown {
			message := "Invalid integration value"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		webhookPath, secret, err := h.services.CommonServices.WebhookService.CreateIntegrationWebhook(ctx, tenant, integration)
		if err != nil {
			message := "Unable to create webhook"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		record := CreateWebhookRecord{
			URL:         fmt.Sprintf("%s%s/%s", baseURL, apiPath, webhookPath),
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

func (h *WebhookHandler) GetActiveWebhooks(baseURL, apiPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetActiveWebhooks", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		integration := c.Query("integration")

		var webhooks *[]postgres_entity.Webhooks
		var err error

		switch {
		case integration == "":
			webhooks, err = h.services.Repositories.PostgresRepositories.WebhooksRepository.FindAll(ctx)
			if err != nil {
				err = fmt.Errorf("Unable to lookup active webhooks for %s: %v", tenant, err)
				tracing.TraceErr(span, err)
				h.responseHandler.HandleError(c, http.StatusNotFound, nil)
				return
			}
		default:
			webhook, err := h.services.CommonServices.WebhookService.GetWebhookForIntegration(ctx, enum.DecodeSource(integration))
			if err != nil || webhook == nil {
				err = fmt.Errorf("Unable to lookup active webhooks for %s: %v", tenant, err)
				tracing.TraceErr(span, err)
				h.responseHandler.HandleError(c, http.StatusNotFound, nil)
				return
			}
			webhookArray := []postgres_entity.Webhooks{*webhook}
			webhooks = &webhookArray
		}

		tracing.LogObjectAsJson(span, "webhooks", webhooks)

		if len(*webhooks) == 0 {
			h.responseHandler.HandleSuccess(c, NoActiveWebhooks{
				Message: "No active webhooks",
			})
			return
		}

		results := make([]ActiveWebhookRecord, len(*webhooks))

		for i, webhook := range *webhooks {
			results[i] = ActiveWebhookRecord{
				URL:         fmt.Sprintf("%s%s/%s", baseURL, apiPath, webhook.WebhookPath),
				Integration: webhook.Integration,
				CreatedAt:   webhook.CreatedAt,
				Active:      webhook.Enabled,
			}
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

func (h *WebhookHandler) RotateWebhook(baseURL, apiPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RotateWebhook", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}
		tenantHash := c.Param("tenantHash")
		validTenant, err := h.services.CommonServices.WebhookService.ValidateTenantId(ctx, tenant, tenantHash)
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
		webhookPath := h.extractWebhookPath(c.Request.URL.Path, apiPath)
		integration, err := h.services.CommonServices.WebhookService.GetIntegrationFromWebhookPath(ctx, tenant, webhookPath)
		if err != nil || integration == commonEnum.SourceUnknown {
			message := "Unable to identify webhook"
			h.responseHandler.HandleError(c, http.StatusNotFound, &message)
			return
		}

		// Call create to rotate webhook as it will automatically handle rotation
		webhookPath, secret, err := h.services.CommonServices.WebhookService.CreateIntegrationWebhook(ctx, tenant, integration)
		if err != nil {
			message := "Unable to rotate webhook"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		record := CreateWebhookRecord{
			URL:         fmt.Sprintf("%s%s/%s", baseURL, apiPath, webhookPath),
			Integration: integration.String(),
			Secret:      secret,
		}

		h.responseHandler.HandleSuccess(c, CreateWebhookResponse{
			Hook: record,
		})
	}
}

func (h *WebhookHandler) DeactivateWebhook(baseURL, apiPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "DeactivateWebhook", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}
		tenantHash := c.Param("tenantHash")
		validTenant, err := h.services.CommonServices.WebhookService.ValidateTenantId(ctx, tenant, tenantHash)
		if err != nil || !validTenant {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		// Call to deactivate webhook
		webhookPath := h.extractWebhookPath(c.Request.URL.Path, apiPath)
		ok, err := h.services.CommonServices.WebhookService.DeactivateWebhookByPath(ctx, webhookPath)
		if err != nil {
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}
		if !ok {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
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

		span.LogFields(log.String("param.tenantHash", c.Param("tenantHash")))

		tenant, err := h.services.Repositories.PostgresRepositories.TenantRepository.GetTenantByHashId(ctx, c.Param("tenantHash"))
		if err != nil {
			err := errors.Wrap(err, "Unable to identify tenant")
			tracing.TraceErr(span, err)
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
		span.LogKV("result.integration", integration.String())

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
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}
	}
}

func (h *WebhookHandler) extractWebhookPath(s, apiRootPath string) string {
	webhookPath := strings.TrimSuffix(s, "/rotate")
	webhookPath = strings.TrimPrefix(webhookPath, apiRootPath)
	webhookPath = strings.Trim(webhookPath, "/")
	return webhookPath
}
