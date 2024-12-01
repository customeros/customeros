package flows

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	comserv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type CreateWebhookRequest struct {
	Integration string `json:"integration"`
}

type CreateWebhookRecord struct {
	URL         string `json:"url"`
	Integration string `json:"integration"`
	Secret      string `json:"secret"`
}

type CreateWebhookResponse struct {
	rest.BaseResponse
	Hook CreateWebhookRecord `json:"hook"`
}

func CreateWebhook(services *service.Services, baseURL, flowsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateWebhook", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrInvalidAPIKey)
			return
		}

		var req CreateWebhookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Missing parameter: integration"))
			return
		}

		integration, err := services.WebhookService.GetIntegration(strings.ToLower(req.Integration))
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("please provide a valid integration value"))
		}

		webhookPath, secret, err := services.WebhookService.CreateIntegrationWebhook(ctx, tenant, comserv.Integration(integration))
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInvalidAPIKey.WithMessage("Unable to create webhook"))
			return
		}

		record := CreateWebhookRecord{
			URL:         fmt.Sprintf("%s%s/%s", baseURL, flowsPath, webhookPath),
			Integration: integration.String(),
			Secret:      secret,
		}

		c.JSON(http.StatusCreated, CreateWebhookResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Hook:         record,
		})
	}
}

type NoActiveWebhooks struct {
	rest.BaseResponse
	Message string `json:"message"`
}

type OneActiveWebhook struct {
	rest.BaseResponse
	Hook ActiveWebhookRecord `json:"hook"`
}

type ActiveWebhooksResponse struct {
	rest.BaseResponse
	Hooks []ActiveWebhookRecord `json:"hooks"`
}

type ActiveWebhookRecord struct {
	URL         string    `json:"url"`
	Integration string    `json:"integration"`
	CreatedAt   time.Time `json:"createdAt"`
	Active      bool      `json:"active"`
}

func GetActiveWebhooks(services *service.Services, baseURL, flowsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetActiveWebhooks", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrInvalidAPIKey)
			return
		}

		count, webhooks, err := services.CommonServices.PostgresRepositories.FlowWebhooksRepository.FindAllActiveWebhooks(ctx, tenant)
		if err != nil {
			err = fmt.Errorf("Unable to lookup active webhooks for %s: %v", tenant, err)
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer)
			return
		}

		if count == 0 {
			c.JSON(http.StatusOK, NoActiveWebhooks{
				BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
				Message:      "No active webhooks",
			})
			return
		}

		results := make([]ActiveWebhookRecord, count)

		for _, webhook := range webhooks {
			record := ActiveWebhookRecord{
				URL:         fmt.Sprintf("%s%s/%s", baseURL, flowsPath, webhook.WebhookPath),
				Integration: webhook.Integration,
				CreatedAt:   webhook.CreatedAt,
				Active:      webhook.Enabled,
			}
			results = append(results, record)
		}

		if count == 1 {
			c.JSON(http.StatusOK, OneActiveWebhook{
				BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
				Hook:         results[0],
			})
			return
		}

		c.JSON(http.StatusOK, ActiveWebhooksResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Hooks:        results,
		})
	}
}

func RotateWebhook(services *service.Services, baseURL, flowsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RotateWebhook", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		// Validate tenant owns webhook
		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrInvalidAPIKey)
			return
		}
		tenantId := c.Param("tenantId")
		validTenant, err := services.WebhookService.ValidateTenantId(ctx, tenant, tenantId)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Unable to verify webhook ownership"))
			return
		}
		if !validTenant {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
			return
		}

		// Lookup integration
		integration, err := services.WebhookService.GetIntegrationFromWebhookPath(ctx, tenant, strings.TrimSuffix(c.Request.URL.Path, "/rotate"))
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Unable to identify webhook"))
			return
		}

		// Call create to rotate webhook as it will automatically handle rotation
		webhookPath, secret, err := services.WebhookService.CreateIntegrationWebhook(ctx, tenant, integration)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Unable to rotate webhook"))
			return
		}

		record := CreateWebhookRecord{
			URL:         fmt.Sprintf("%s%s/%s", baseURL, flowsPath, webhookPath),
			Integration: integration.String(),
			Secret:      secret,
		}

		c.JSON(http.StatusCreated, CreateWebhookResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Hook:         record,
		})
	}
}

func DeactivateWebhook(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RotateWebhook", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		// Validate tenant owns webhook
		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrInvalidAPIKey)
			return
		}
		tenantId := c.Param("tenantId")
		validTenant, err := services.WebhookService.ValidateTenantId(ctx, tenant, tenantId)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Unable to verify webhook ownership"))
			return
		}
		if !validTenant {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
			return
		}

		// Call to deactivate webhook
		deactErr := services.WebhookService.DeactivateWebhook(ctx, strings.TrimSuffix(c.Request.URL.Path, "/rotate"))
		if deactErr != nil {
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer)
			return
		}

		c.JSON(http.StatusOK, NoActiveWebhooks{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Message:      "Webhook successfully deactivated",
		})
	}
}

func HandleWebhook(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Webhooks", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant, err := s.CommonServices.PostgresRepositories.TenantRepository.GetTenant(ctx, c.Param("tenantId"))
		if err != nil {
			err := errors.Wrap(err, "Unable to identify tenant")
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
			return
		}

		integration, err := s.WebhookService.GetIntegrationFromWebhookPath(ctx, tenant, c.Request.URL.Path)
		if err != nil {
			rest.SendError(c, span, http.StatusNotFound, rest.ErrNotFound.WithMessage("Webhook not found"))
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
		case comserv.IntegrationCalCom:
			CalDotCom(&httpContext)
		// todo
		case comserv.IntegrationFathom:
			FathomZapier(&httpContext)
		case comserv.IntegrationGrain:
			GrainZapier(&httpContext)
		case comserv.IntegrationPostmark:
			PostmarkInboundEmail(&httpContext)
		// todo
		default:
			rest.SendError(c, span, http.StatusNotFound, rest.ErrNotFound)
			return
		}
	}
}
