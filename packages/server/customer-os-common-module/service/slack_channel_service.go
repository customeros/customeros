package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"time"
)

type SlackChannelService interface {
	StoreSlackChannel(tenant, channel string, organizationId *string) error
}

type slackChannelService struct {
	repositories *repository.Repositories
}

func NewSlackChannelService(repositories *repository.Repositories) SlackChannelService {
	return &slackChannelService{
		repositories: repositories,
	}
}

func (u *slackChannelService) StoreSlackChannel(tenant, channelId string, organizationId *string) error {
	existing, err := u.repositories.SlackChannelRepository.GetSlackChannel(tenant, channelId)
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
		}
		return u.repositories.SlackChannelRepository.CreateSlackChannel(&slackChannel)
	}

	if existing != nil && organizationId != nil {
		return u.repositories.SlackChannelRepository.UpdateSlackChannel(existing.ID, *organizationId)
	}

	return nil
}
