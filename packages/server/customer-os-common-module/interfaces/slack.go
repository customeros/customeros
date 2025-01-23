package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type SlackService interface {
	GetSlackChannels(ctx context.Context, tenant string) ([]*postgres_entity.SlackChannel, error)
	GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) ([]*postgres_entity.SlackChannel, int64, error)
	StoreSlackChannel(ctx context.Context, tenant, source, channelId, channelName string, organizationId *string) error
	SendMessageFromBot(ctx context.Context, channel, blocks string) error
}
