package webhook

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"strings"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
)

type webhookService struct {
	log                  logger.Logger
	postgresRepositories *postgres_repository.Repositories
}

func NewWebhookService(log logger.Logger, postgresRepositories *postgres_repository.Repositories) interfaces.WebhookService {
	return &webhookService{
		log:                  log,
		postgresRepositories: postgresRepositories,
	}
}

func (w *webhookService) GetIntegration(s string) enum.Source {
	integration := enum.DecodeSource(s)
	switch integration {
	case enum.SourceGrain, enum.SourceFathom:
		return integration
	default:
		return enum.SourceUnknown
	}
}

func (w *webhookService) ValidateTenantId(ctx context.Context, tenant, tenantId string) (bool, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebhookService.ValidateTenantId")
	defer spans.Finish()
	spans.LogKV("tenant", tenant)
	spans.LogKV("tenantId", tenantId)

	tenantFromDb, err := w.postgresRepositories.TenantRepository.GetTenantByHashId(ctx, tenantId)
	if err != nil {
		err = fmt.Errorf("Unable to lookup tenant hashId for %s: %v", tenant, err)
		spans.TraceError(err)
	}

	return tenantFromDb == tenant, nil
}

func (w *webhookService) GetIntegrationFromWebhookPath(ctx context.Context, tenant, webhookPath string) (enum.Source, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebhookService.GetIntegrationFromWebhookPath")
	defer spans.Finish()
	spans.LogKV("webhookPath", webhookPath)

	path := strings.Trim(webhookPath, "/")
	webhook, err := w.postgresRepositories.WebhooksRepository.Find(ctx, postgres_entity.Webhooks{
		Tenant:      tenant,
		WebhookPath: path,
	})
	if err != nil {
		err = fmt.Errorf("Unable to lookup webhook path: %v", err)
		spans.TraceError(err)
		return enum.SourceUnknown, err
	}

	if webhook == nil {
		return enum.SourceUnknown, nil
	}

	if !webhook.Enabled {
		err = fmt.Errorf("Webhook is disabled: %v", err)
		spans.TraceError(err)
		return enum.SourceUnknown, err
	}

	return enum.DecodeSource(webhook.Integration), nil
}

func (w *webhookService) GetWebhookForIntegration(ctx context.Context, integration enum.Source) (*postgres_entity.Webhooks, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebhookService.GetWebhookForIntegration")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		return nil, coserrors.ErrTenantNotSet
	}

	webhook, err := w.postgresRepositories.WebhooksRepository.Find(ctx, postgres_entity.Webhooks{
		Tenant:      tenant,
		Enabled:     true,
		Integration: integration.String(),
	})

	return webhook, err
}
