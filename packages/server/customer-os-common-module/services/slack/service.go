package slack

import (
	"context"
	"time"

	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type slackService struct {
	postgres *repository.Repositories
}

func NewSlackService(postgres *repository.Repositories) interfaces.SlackService {
	return &slackService{
		postgres: postgres,
	}
}

func (s *slackService) GetSlackChannels(ctx context.Context, tenant string) ([]*postgresEntity.SlackChannel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.GetSlackChannels")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, "service")

	nodes, err := s.postgres.SlackChannelRepository.GetSlackChannels(ctx, tenant)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (s *slackService) GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) ([]*postgresEntity.SlackChannel, int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.GetSlackChannels")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, "service")

	channels, totalCount, err := s.postgres.SlackChannelRepository.GetPaginatedSlackChannels(ctx, tenant, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return channels, totalCount, nil
}

func (s *slackService) StoreSlackChannel(ctx context.Context, tenant, source, channelId, channelName string, organizationId *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackService.StoreSlackChannel")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, "service")
	span.LogFields(log.String("channelId", channelId))
	span.LogFields(log.String("channelName", channelName))
	span.LogFields(log.String("organizationId", utils.IfNotNilString(organizationId)))

	existing, err := s.postgres.SlackChannelRepository.GetSlackChannel(ctx, tenant, channelId)
	if err != nil {
		return err
	}

	if existing == nil {
		now := time.Now()
		slackChannel := postgresEntity.SlackChannel{
			CreatedAt:      now,
			UpdatedAt:      now,
			TenantName:     tenant,
			ChannelId:      channelId,
			ChannelName:    channelName,
			OrganizationId: organizationId,
			Source:         source,
		}
		return s.postgres.SlackChannelRepository.CreateSlackChannel(ctx, &slackChannel)
	}
	if existing != nil {
		if organizationId != nil {
			return s.postgres.SlackChannelRepository.UpdateSlackChannelOrganization(ctx, existing.ID, *organizationId)
		} else if channelName != "" {
			return s.postgres.SlackChannelRepository.UpdateSlackChannelName(ctx, existing.ID, channelName)
		}
	}

	return nil
}
