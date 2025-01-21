package cosapi_interfaces

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type SlackService interface {
	GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) (*utils.Pagination, error)
	GetSlackSettings(ctx context.Context, tenant string) (*SlackSettingsResponse, error)
	ListSlackChannelsWithBot(ctx context.Context, tenant string) ([]SlackChannelResponse, error)
	JoinSlackChannelsWithBot(ctx context.Context, tenant, channelId string) error
	LeaveSlackChannelsWithBot(ctx context.Context, tenant, channelId string) error
}

type SlackSettingsResponse struct {
	SlackEnabled bool `json:"slackEnabled"`
}

type SlackChannelResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
