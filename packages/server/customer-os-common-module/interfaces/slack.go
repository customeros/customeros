package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type SlackService interface {
	GetSlackChannels(ctx context.Context, tenant string) ([]*postgres_entity.SlackChannel, error)
	GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) (*utils.Pagination, error)
	StoreSlackChannel(ctx context.Context, tenant, source, channelId, channelName string, organizationId *string) error
	SendMessageFromBot(ctx context.Context, channel, blocks string) error
	GetSlackSettings(ctx context.Context, tenant string) (*SlackSettingsResponse, error)
	ListSlackChannelsWithBot(ctx context.Context) ([]SlackChannelResponse, error)
	JoinSlackChannelsWithBot(ctx context.Context, channelId string) error
	LeaveSlackChannelsWithBot(ctx context.Context, channelId string) error
}

type SlackSettingsResponse struct {
	SlackEnabled bool `json:"slackEnabled"`
}

type SlackChannelResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
