package service

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
)

type WebhookService interface {
	GetIntegration(s string) (Integration, error)
	CreateIntegrationWebhook(ctx context.Context, tenant string, integration Integration) (webhookUrl string, secret string, err error)
}

type Integration string

// Add all supported integrations here, and also in validIntegrations below
const (
	IntegrationCalCom   Integration = "calcom"
	IntegrationFathom   Integration = "fathom"
	IntegrationGrain    Integration = "grain"
	IntegrationPostmark Integration = "postmark"
)

var validIntegrations = func() map[string]Integration {
	integrations := []Integration{
		IntegrationCalCom,
		IntegrationFathom,
		IntegrationGrain,
		IntegrationPostmark,
	}

	m := make(map[string]Integration)
	for _, i := range integrations {
		m[string(i)] = i
	}
	return m
}()

func (i Integration) String() string {
	return string(i)
}

func (i Integration) IntegrationID(rotationCount int64) string {
	// Add rotation count to string being hashed
	input := fmt.Sprintf("%s:%d", i.String(), rotationCount)
	return utils.GenerateHashId(input, 12)
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

func (w *webhookService) GetIntegration(s string) (Integration, error) {
	if integration, ok := validIntegrations[s]; ok {
		return integration, nil
	}
	return "", fmt.Errorf("invalid integration type: %s", s)
}

func (w *webhookService) CreateIntegrationWebhook(ctx context.Context, tenant string, integration Integration) (string, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhookService.CreateWebhook")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))
	span.LogFields(log.String("integration", integration.String()))

	var newWebhook entity.FlowWebhooks

	tenantHash, err := w.repositories.PostgresRepositories.TenantRepository.GetHashID(ctx, tenant)
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

	integrationHash := integration.IntegrationID(int64(webhook.RotationCount + 1))
	secret, err := utils.GenerateSecret()
	if err != nil {
		err = fmt.Errorf("Unable to generate webhook secret: %v", err)
	}

	newWebhook.TenantName = tenant
	newWebhook.WebhookPath = fmt.Sprintf("%s/i/%s", tenantHash, integrationHash)
	newWebhook.Integration = integration.String()
	newWebhook.Secret = secret
	newWebhook.RotationCount = 1

	// disable existing webhook for tenant/integration if exists
	if count != 0 {
		err := w.repositories.PostgresRepositories.FlowWebhooksRepository.DisableWebhook(ctx, webhook.WebhookPath)
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

	createErr := w.repositories.PostgresRepositories.FlowWebhooksRepository.CreateWebhook(&newWebhook)
	if createErr != nil {
		err = fmt.Errorf("Unable to create webhook: %v", err)
		tracing.TraceErr(span, err)
	}

	return newWebhook.WebhookPath, newWebhook.Secret, nil
}
