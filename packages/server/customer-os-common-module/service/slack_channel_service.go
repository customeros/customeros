package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type SlackChannelService interface {
	GetSlackChannels(ctx context.Context, tenant string) ([]*entity.SlackChannel, error)
	GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) ([]*entity.SlackChannel, int64, error)

	StoreSlackChannel(ctx context.Context, tenant, source, channel string, organizationId *string) error
}

type slackChannelService struct {
	repositories *repository.Repositories
}

func NewSlackChannelService(repositories *repository.Repositories) SlackChannelService {
	return &slackChannelService{
		repositories: repositories,
	}
}

func (s *slackChannelService) GetSlackChannels(ctx context.Context, tenant string) ([]*entity.SlackChannel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetSlackChannels.GetSlackChannels")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "service")
	span.LogFields(log.String("tenant", tenant))

	nodes, err := s.repositories.SlackChannelRepository.GetSlackChannels(tenant)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (s *slackChannelService) GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) ([]*entity.SlackChannel, int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetSlackChannels.GetSlackChannels")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "service")
	span.LogFields(log.String("tenant", tenant))

	channels, totalCount, err := s.repositories.SlackChannelRepository.GetPaginatedSlackChannels(tenant, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return channels, totalCount, nil
}

func (s *slackChannelService) StoreSlackChannel(ctx context.Context, tenant, source, channelId string, organizationId *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetSlackChannels.StoreSlackChannel")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "service")
	span.LogFields(log.String("tenant", tenant))
	span.LogFields(log.String("channelId", channelId))
	span.LogFields(log.String("organizationId", utils.IfNotNilString(organizationId)))

	existing, err := s.repositories.SlackChannelRepository.GetSlackChannel(tenant, channelId)
	if err != nil {
		return err
	}

	if existing == nil {
		now := time.Now()
		slackChannel := entity.SlackChannel{
			CreatedAt:      now,
			UpdatedAt:      now,
			TenantName:     tenant,
			ChannelId:      channelId,
			OrganizationId: organizationId,
			Source:         source,
		}
		return s.repositories.SlackChannelRepository.CreateSlackChannel(&slackChannel)
	}

	if existing != nil && organizationId != nil {
		return s.repositories.SlackChannelRepository.UpdateSlackChannel(existing.ID, *organizationId)
	}

	return nil
}
