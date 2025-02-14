package webhook

import (
	"context"
	"fmt"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

func (w *webhookService) CreateIntegrationWebhook(ctx context.Context, tenant string, integration enum.Source) (string, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhookService.CreateIntegrationWebhook")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))
	span.LogFields(log.String("integration", integration.String()))

	rotationCount, err := w.postgresRepositories.WebhooksRepository.FindLastRotationCount(ctx, integration)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", "", err
	}

	newWebhook, err := w.buildWebhook(ctx, integration, rotationCount)
	if err != nil || newWebhook == nil {
		tracing.TraceErr(span, err)
		return "", "", err
	}

	result, err := w.postgresRepositories.WebhooksRepository.Create(ctx, *newWebhook)
	if err != nil {
		err = fmt.Errorf("Unable to create webhook: %v", err)
		tracing.TraceErr(span, err)
	}

	return result.WebhookPath, result.Secret, nil
}

func (w *webhookService) buildWebhook(ctx context.Context, integration enum.Source, rotationCount int) (*postgres_entity.Webhooks, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)
	rotationCount++

	tenantHash, err := w.postgresRepositories.TenantRepository.GetHashID(ctx, tenant)
	span.LogFields(log.String("tenantHash", tenantHash))
	if err != nil {
		err = fmt.Errorf("Unable to get HashID for tenant %s: %v", tenant, err)
		tracing.TraceErr(span, err)
		return nil, err
	}

	integrationHash := integration.IntegrationID(rotationCount)
	secret, err := utils.GenerateSecret()
	if err != nil {
		err = fmt.Errorf("Unable to generate webhook secret: %v", err)
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &webhook, nil
}
