package cosapi_interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type SlackService interface {
	GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) (*utils.Pagination, error)
	GetSlackSettings(ctx context.Context, tenant string) (*SlackSettingsResponse, error)
}

type SlackSettingsResponse struct {
	SlackEnabled bool `json:"slackEnabled"`
}
