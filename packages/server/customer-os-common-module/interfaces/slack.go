package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type SlackService interface {
	GetSlackChannels(ctx context.Context, tenant string) ([]*entity.SlackChannel, error)
	GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) ([]*entity.SlackChannel, int64, error)
	StoreSlackChannel(ctx context.Context, tenant, source, channelId, channelName string, organizationId *string) error
}
