package webhook

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

func (w *webhookService) CreateIntegrationWebhook(ctx context.Context, tenant string, integration enum.Source) (string, string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebhookService.CreateIntegrationWebhook")
	defer spans.Finish()
	spans.LogKV("integration", integration.String())

	rotationCount, err := w.postgresRepositories.WebhooksRepository.FindLastRotationCount(ctx, integration)
	if err != nil {
		spans.TraceError(err)
		return "", "", err
	}

	newWebhook, err := w.buildWebhook(ctx, integration, rotationCount)
	if err != nil || newWebhook == nil {
		spans.TraceError(err)
		return "", "", err
	}

	result, err := w.postgresRepositories.WebhooksRepository.Create(ctx, *newWebhook)
	if err != nil {
		err = fmt.Errorf("Unable to create webhook: %v", err)
		spans.TraceError(err)
	}

	return result.WebhookPath, result.Secret, nil
}

func (w *webhookService) buildWebhook(ctx context.Context, integration enum.Source, rotationCount int) (*postgres_entity.Webhooks, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)
	rotationCount++

	tenantHash, err := w.postgresRepositories.TenantRepository.GetHashID(ctx, tenant)
	spans.LogKV("tenantHash", tenantHash)
	if err != nil {
		err = fmt.Errorf("Unable to get HashID for tenant %s: %v", tenant, err)
		spans.TraceError(err)
		return nil, err
	}

	integrationHash := integration.IntegrationID(rotationCount)
	secret, err := utils.GenerateSecret()
	if err != nil {
		err = fmt.Errorf("Unable to generate webhook secret: %v", err)
		spans.TraceError(err)
		return nil, err
	}

	webhook := postgres_entity.Webhooks{
		Tenant:        common.GetTenantFromContext(ctx),
		WebhookPath:   fmt.Sprintf("%s/i/%s", tenantHash, integrationHash),
		Integration:   integration.String(),
		Secret:        secret,
		RotationCount: rotationCount,
	}

	err = webhook.Validate()
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &webhook, nil
}
