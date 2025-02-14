package webhook

import (
	"context"
	"fmt"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

func (w *webhookService) DeactivateWebhookByPath(ctx context.Context, webhookPath string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhookService.DeactivateWebhookByPath")
	defer span.Finish()
	span.LogKV("webhookPath", webhookPath)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		tracing.TraceErr(span, coserrors.ErrTenantNotSet)
		return false, coserrors.ErrTenantNotSet
	}

	query := postgres_entity.Webhooks{
		Tenant:      tenant,
		WebhookPath: webhookPath,
	}

	ok, err := w.postgresRepositories.WebhooksRepository.Deactivate(ctx, query)
	if err != nil {
		err = fmt.Errorf("Unable to deactivate webhook %s: %v", webhookPath, err)
		tracing.TraceErr(span, err)
		return false, err
	}

	return ok, nil
}
