package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type WebhookService interface {
	GetIntegration(s string) (enum.Source, error)
	GetWebhookForIntegration(ctx context.Context, integration enum.Source) (*postgres_entity.Webhooks, error)
	CreateIntegrationWebhook(ctx context.Context, tenant string, integration enum.Source) (webhookUrl string, secret string, err error)
	ValidateTenantId(ctx context.Context, tenant, tenantId string) (bool, error)
	GetIntegrationFromWebhookPath(ctx context.Context, tenant, webhookPath string) (enum.Source, error)
	DeactivateWebhook(ctx context.Context, webhookPath string) error
}
