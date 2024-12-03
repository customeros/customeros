package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
)

type WebhookService interface {
	GetIntegration(s string) (enum.ExternalSystemId, error)
	CreateIntegrationWebhook(ctx context.Context, tenant string, integration enum.ExternalSystemId) (webhookUrl string, secret string, err error)
	ValidateTenantId(ctx context.Context, tenant, tenantId string) (bool, error)
	GetIntegrationFromWebhookPath(ctx context.Context, tenant, webhookPath string) (enum.ExternalSystemId, error)
	DeactivateWebhook(ctx context.Context, webhookPath string) error
}

type webhookService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewWebhookService(log logger.Logger, repositories *repository.Repositories, services *Services) WebhookService {
	return &webhookService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (w *webhookService) GetIntegration(s string) (enum.ExternalSystemId, error) {
	integration := enum.DecodeExternalSystemId(s)
	if integration == "" {
		return "", fmt.Errorf("Invalid integration %s", s)
	}
	return integration, nil
}

func (w *webhookService) ValidateTenantId(ctx context.Context, tenant, tenantId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhookService.ValidateTenantId")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))
	span.LogFields(log.String("tenantId", tenantId))

	tenantFromDb, err := w.services.Repositories.PostgresRepositories.TenantRepository.GetTenant(ctx, tenantId)
	if err != nil {
		err = fmt.Errorf("Unable to lookup tenant hashId for %s: %v", tenant, err)
		tracing.TraceErr(span, err)
	}

	return tenantFromDb == tenant, nil
}

func (w *webhookService) GetIntegrationFromWebhookPath(ctx context.Context, tenant, webhookPath string) (enum.ExternalSystemId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhookService.GetIntegrationFromWebhookPath")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))
	span.LogFields(log.String("webhookPath", webhookPath))

	path := strings.TrimPrefix(webhookPath, "/")
	webhook, err := w.services.Repositories.PostgresRepositories.FlowWebhooksRepository.FindWebhookByPath(ctx, tenant, path)
	if err != nil {
		err = fmt.Errorf("Unable to lookup webhook path: %v", err)
		tracing.TraceErr(span, err)
	}

	return enum.DecodeExternalSystemId(webhook.Integration), nil
}

func (w *webhookService) CreateIntegrationWebhook(ctx context.Context, tenant string, integration enum.ExternalSystemId) (string, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhookService.CreateIntegrationWebhook")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))
	span.LogFields(log.String("integration", integration.String()))

	var newWebhook entity.FlowWebhooks

	tenantHash, err := w.repositories.PostgresRepositories.TenantRepository.GetHashID(ctx, tenant)
	span.LogFields(log.String("tenantHash", tenantHash))
	if err != nil {
		err = fmt.Errorf("Unable to get HashID for tenant %s: %v", tenant, err)
		tracing.TraceErr(span, err)
		return "", "", err
	}

	// check if webhook already exists for tenant/integration
	count, webhook, err := w.repositories.PostgresRepositories.FlowWebhooksRepository.FindActiveWebhook(ctx, tenant, integration.String())
	if err != nil {
		err = fmt.Errorf("Unable to check db for existing webhook: %v", err)
		tracing.TraceErr(span, err)
		return "", "", err
	}

	integrationHash := integration.IntegrationID(webhook.RotationCount + 1)
	span.LogFields(log.String("integrationHash", integrationHash))
	secret, err := utils.GenerateSecret()
	span.LogFields(log.String("secret", secret))
	if err != nil {
		err = fmt.Errorf("Unable to generate webhook secret: %v", err)
		tracing.TraceErr(span, err)
	}

	newWebhook.TenantName = tenant
	newWebhook.WebhookPath = fmt.Sprintf("%s/i/%s", tenantHash, integrationHash)
	newWebhook.Integration = integration.String()
	newWebhook.Secret = secret
	newWebhook.RotationCount = 1

	// disable existing webhook for tenant/integration if exists
	if count != 0 {
		err := w.DeactivateWebhook(ctx, webhook.WebhookPath)
		if err != nil {
			err = fmt.Errorf("Unable to deactivate existing webhook for %s and %s: %v", tenant, integration.String(), err)
			tracing.TraceErr(span, err)
			return "", "", err
		}

		newWebhook.RotationCount = webhook.RotationCount + 1
	}

	validationErr := newWebhook.Validate()
	if validationErr != nil {
		tracing.TraceErr(span, validationErr)
		return "", "", err
	}

	createErr := w.repositories.PostgresRepositories.FlowWebhooksRepository.CreateWebhook(ctx, &newWebhook)
	if createErr != nil {
		err = fmt.Errorf("Unable to create webhook: %v", err)
		tracing.TraceErr(span, err)
	}

	return newWebhook.WebhookPath, newWebhook.Secret, nil
}

func (w *webhookService) DeactivateWebhook(ctx context.Context, webhookPath string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhookService.Deactivate")
	defer span.Finish()
	span.LogFields(log.String("webhookPath", webhookPath))

	err := w.repositories.PostgresRepositories.FlowWebhooksRepository.DisableWebhook(ctx, webhookPath)
	if err != nil {
		err = fmt.Errorf("Unable to deactivate webhook %s: %v", webhookPath, err)
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}
