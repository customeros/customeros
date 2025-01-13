package flows

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	integrations "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/flows_integrations"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
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
	enum.BaseResponse
	Hook CreateWebhookRecord `json:"hook"`
}

func CreateWebhook(s *cosapi_services.Services, baseURL, flowsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateWebhook", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := handlers.ValidateTenant(c, ctx, span)
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		var req CreateWebhookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Missing parameter: integration"))
			return
		}

		integration, err := s.WebhookService.GetIntegration(strings.ToLower(req.Integration))
		if err != nil {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("please provide a valid integration value"))
			return
		}

		webhookPath, secret, err := s.WebhookService.CreateIntegrationWebhook(ctx, tenant, integration)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Unable to create webhook"))
			return
		}

		record := CreateWebhookRecord{
			URL:         fmt.Sprintf("%s%s/%s", baseURL, flowsPath, webhookPath),
			Integration: integration.String(),
			Secret:      secret,
		}

		c.JSON(http.StatusCreated, CreateWebhookResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Hook:         record,
		})
	}
}

type NoActiveWebhooks struct {
	enum.BaseResponse
	Message string `json:"message"`
}

type OneActiveWebhook struct {
	enum.BaseResponse
	Hook ActiveWebhookRecord `json:"hook"`
}

type ActiveWebhooksResponse struct {
	enum.BaseResponse
	Hooks []ActiveWebhookRecord `json:"hooks"`
}

type ActiveWebhookRecord struct {
	URL         string    `json:"url"`
	Integration string    `json:"integration"`
	CreatedAt   time.Time `json:"createdAt"`
	Active      bool      `json:"active"`
}

func GetActiveWebhooks(s *cosapi_services.Services, baseURL, flowsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetActiveWebhooks", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := handlers.ValidateTenant(c, ctx, span)
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		webhooks, err := s.Repositories.PostgresRepositories.WebhooksRepository.FindAll(ctx)
		if err != nil {
			err = fmt.Errorf("Unable to lookup active webhooks for %s: %v", tenant, err)
			tracing.TraceErr(span, err)
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}

		if len(*webhooks) == 0 {
			c.JSON(http.StatusOK, NoActiveWebhooks{
				BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
				Message:      "No active webhooks",
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
			c.JSON(http.StatusOK, OneActiveWebhook{
				BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
				Hook:         results[0],
			})
			return
		}

		c.JSON(http.StatusOK, ActiveWebhooksResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Hooks:        results,
		})
	}
}

func RotateWebhook(s *cosapi_services.Services, baseURL, flowsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RotateWebhook", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		// Validate tenant owns webhook
		tenant := handlers.ValidateTenant(c, ctx, span)
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}
		tenantId := c.Param("tenantId")
		validTenant, err := s.WebhookService.ValidateTenantId(ctx, tenant, tenantId)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Unable to verify webhook ownership"))
			return
		}
		if !validTenant {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			return
		}

		// Lookup integration
		webhookPath := strings.TrimSuffix(c.Request.URL.Path, "/rotate")
		webhookPath = strings.TrimPrefix(webhookPath, flowsPath)
		integration, err := s.WebhookService.GetIntegrationFromWebhookPath(ctx, tenant, webhookPath)
		if err != nil || integration == commonEnum.SourceUnknown {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("Unable to identify webhook"))
			return
		}

		// Call create to rotate webhook as it will automatically handle rotation
		webhookPath, secret, err := s.WebhookService.CreateIntegrationWebhook(ctx, tenant, integration)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Unable to rotate webhook"))
			return
		}

		record := CreateWebhookRecord{
			URL:         fmt.Sprintf("%s%s/%s", baseURL, flowsPath, webhookPath),
			Integration: integration.String(),
			Secret:      secret,
		}

		c.JSON(http.StatusCreated, CreateWebhookResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Hook:         record,
		})
	}
}

func DeactivateWebhook(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RotateWebhook", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		// Validate tenant owns webhook
		tenant := handlers.ValidateTenant(c, ctx, span)
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}
		tenantId := c.Param("tenantId")
		validTenant, err := s.WebhookService.ValidateTenantId(ctx, tenant, tenantId)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Unable to verify webhook ownership"))
			return
		}
		if !validTenant {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			return
		}

		// Call to deactivate webhook
		deactErr := s.WebhookService.DeactivateWebhook(ctx, strings.TrimSuffix(c.Request.URL.Path, "/rotate"))
		if deactErr != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}

		c.JSON(http.StatusOK, NoActiveWebhooks{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Message:      "Webhook successfully deactivated",
		})
	}
}

func HandleWebhook(s *cosapi_services.Services, flowsPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Webhooks", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant, err := s.Repositories.PostgresRepositories.TenantRepository.GetTenant(ctx, c.Param("tenantId"))
		if err != nil {
			err := errors.Wrap(err, "Unable to identify tenant")
			tracing.TraceErr(span, err)
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			return
		}

		webhookPath := strings.TrimPrefix(c.Request.URL.Path, flowsPath)
		integration, err := s.WebhookService.GetIntegrationFromWebhookPath(ctx, tenant, webhookPath)
		if err != nil {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("Webhook not found"))
			return
		}

		switch integration {
		case commonEnum.SourceCalCom:
			integrations.CalDotCom(c, s)
		// todo
		case commonEnum.SourceFathom:
			integrations.FathomZapier(c, s)
		case commonEnum.SourceGrain:
			integrations.GrainZapier(c, s)
		case commonEnum.SourcePostmark:
			integrations.PostmarkInboundEmail(c, s)
		// todo
		default:
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound)
			return
		}
	}
}
