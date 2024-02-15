package repository

import (
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"gorm.io/gorm"
)

type slackChannelRepo struct {
	db *gorm.DB
}

type SlackChannelRepository interface {
	GetSlackChannel(tenant, channelId string) (*entity.SlackChannel, error)
	GetSlackChannels(tenant string) ([]*entity.SlackChannel, error)
	CreateSlackChannel(entity *entity.SlackChannel) error
	UpdateSlackChannel(entityId uint64, organizationId string) error
}

func NewSlackChannelRepository(db *gorm.DB) SlackChannelRepository {
	return &slackChannelRepo{db: db}
}

func (r *slackChannelRepo) GetSlackChannel(tenant, channelId string) (*entity.SlackChannel, error) {
	var entities []entity.SlackChannel
	err := r.db.
		Where("tenant_name = ?", tenant).
		Where("channel_id = ?", channelId).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	if len(entities) == 0 {
		return nil, nil
	}
	if len(entities) > 1 {
		return nil, errors.New("multiple slack channels found")
	}

	return &entities[0], nil
}

func (r *slackChannelRepo) GetSlackChannels(tenant string) ([]*entity.SlackChannel, error) {
	var entities []*entity.SlackChannel
	err := r.db.
		Where("tenant_name = ?", tenant).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *slackChannelRepo) CreateSlackChannel(entity *entity.SlackChannel) error {
	err := r.db.Create(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *slackChannelRepo) UpdateSlackChannel(entityId uint64, organizationId string) error {
	return r.db.Model(&entity.SlackChannel{}).Where("id = ?", entityId).Update("organization_id", organizationId).Error
}
