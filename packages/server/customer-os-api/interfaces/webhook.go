package cosapi_interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

type WebhookService interface {
	GetIntegration(s string) (enum.Source, error)
	CreateIntegrationWebhook(ctx context.Context, tenant string, integration enum.Source) (webhookUrl string, secret string, err error)
	ValidateTenantId(ctx context.Context, tenant, tenantId string) (bool, error)
	GetIntegrationFromWebhookPath(ctx context.Context, tenant, webhookPath string) (enum.Source, error)
	DeactivateWebhook(ctx context.Context, webhookPath string) error
}
