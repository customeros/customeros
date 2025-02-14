package webhook

import (
	"context"
	"fmt"
	"strings"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhookService.ValidateTenantId")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))
	span.LogFields(log.String("tenantId", tenantId))

	tenantFromDb, err := w.postgresRepositories.TenantRepository.GetTenant(ctx, tenantId)
	if err != nil {
		err = fmt.Errorf("Unable to lookup tenant hashId for %s: %v", tenant, err)
		tracing.TraceErr(span, err)
	}

	return tenantFromDb == tenant, nil
}

func (w *webhookService) GetIntegrationFromWebhookPath(ctx context.Context, tenant, webhookPath string) (enum.Source, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhookService.GetIntegrationFromWebhookPath")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("webhookPath", webhookPath))

	path := strings.Trim(webhookPath, "/")
	webhook, err := w.postgresRepositories.WebhooksRepository.Find(ctx, postgres_entity.Webhooks{
		Tenant:      tenant,
		WebhookPath: path,
	})
	if err != nil {
		err = fmt.Errorf("Unable to lookup webhook path: %v", err)
		tracing.TraceErr(span, err)
		return enum.SourceUnknown, err
	}

	if webhook == nil {
		return enum.SourceUnknown, nil
	}

	if !webhook.Enabled {
		err = fmt.Errorf("Webhook is disabled: %v", err)
		tracing.TraceErr(span, err)
		return enum.SourceUnknown, err
	}

	return enum.DecodeSource(webhook.Integration), nil
}

func (w *webhookService) GetWebhookForIntegration(ctx context.Context, integration enum.Source) (*postgres_entity.Webhooks, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhookService.GetWebhookForIntegration")
	defer span.Finish()
	tracing.TagComponentService(span)

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
