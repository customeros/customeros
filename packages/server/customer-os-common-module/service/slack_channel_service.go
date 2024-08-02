package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type SlackChannelService interface {
	GetSlackChannels(ctx context.Context, tenant string) ([]*postgresEntity.SlackChannel, error)
	GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) ([]*postgresEntity.SlackChannel, int64, error)

	StoreSlackChannel(ctx context.Context, tenant, source, channelId, channelName string, organizationId *string) error
}

type slackChannelService struct {
	repositories *postgresRepository.Repositories
}

func NewSlackChannelService(repositories *postgresRepository.Repositories) SlackChannelService {
	return &slackChannelService{
		repositories: repositories,
	}
}

func (s *slackChannelService) GetSlackChannels(ctx context.Context, tenant string) ([]*postgresEntity.SlackChannel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetSlackChannels.GetSlackChannels")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, "service")

	nodes, err := s.repositories.SlackChannelRepository.GetSlackChannels(tenant)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (s *slackChannelService) GetPaginatedSlackChannels(ctx context.Context, tenant string, page, limit int) ([]*postgresEntity.SlackChannel, int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetSlackChannels.GetSlackChannels")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, "service")

	channels, totalCount, err := s.repositories.SlackChannelRepository.GetPaginatedSlackChannels(tenant, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return channels, totalCount, nil
}

func (s *slackChannelService) StoreSlackChannel(ctx context.Context, tenant, source, channelId, channelName string, organizationId *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetSlackChannels.StoreSlackChannel")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, "service")
	span.LogFields(log.String("channelId", channelId))
	span.LogFields(log.String("channelName", channelName))
	span.LogFields(log.String("organizationId", utils.IfNotNilString(organizationId)))

	existing, err := s.repositories.SlackChannelRepository.GetSlackChannel(tenant, channelId)
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
		return s.repositories.SlackChannelRepository.CreateSlackChannel(&slackChannel)
	}
	if existing != nil {
		if organizationId != nil {
			return s.repositories.SlackChannelRepository.UpdateSlackChannelOrganization(existing.ID, *organizationId)
		} else if channelName != "" {
			return s.repositories.SlackChannelRepository.UpdateSlackChannelName(existing.ID, channelName)
		}
	}

	return nil
}
