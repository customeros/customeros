package webhook

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
)

func (w *webhookService) DeactivateWebhookByPath(ctx context.Context, webhookPath string) (bool, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebhookService.DeactivateWebhookByPath")
	defer spans.Finish()
	spans.LogKV("webhookPath", webhookPath)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		spans.TraceError(coserrors.ErrTenantNotSet)
		return false, coserrors.ErrTenantNotSet
	}

	query := postgres_entity.Webhooks{
		Tenant:      tenant,
		WebhookPath: webhookPath,
	}

	ok, err := w.postgresRepositories.WebhooksRepository.Deactivate(ctx, query)
	if err != nil {
		err = fmt.Errorf("Unable to deactivate webhook %s: %v", webhookPath, err)
		spans.TraceError(err)
		return false, err
	}

	return ok, nil
}
