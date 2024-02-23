package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"gorm.io/gorm"
)

type slackChannelRepository struct {
	db *gorm.DB
}

type SlackChannelRepository interface {
	GetSlackChannel(tenant, channelId string) (*entity.SlackChannel, error)
	GetSlackChannels(tenant string) ([]*entity.SlackChannel, error)
	GetPaginatedSlackChannels(tenant string, skip, limit int) ([]*entity.SlackChannel, int64, error)

	CreateSlackChannel(entity *entity.SlackChannel) error
	UpdateSlackChannel(entityId uuid.UUID, organizationId string) error
}

func NewSlackChannelRepository(db *gorm.DB) SlackChannelRepository {
	return &slackChannelRepository{db: db}
}

func (r *slackChannelRepository) GetSlackChannel(tenant, channelId string) (*entity.SlackChannel, error) {
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

func (r *slackChannelRepository) GetSlackChannels(tenant string) ([]*entity.SlackChannel, error) {
	var entities []*entity.SlackChannel
	err := r.db.
		Where("tenant_name = ?", tenant).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *slackChannelRepository) GetPaginatedSlackChannels(tenant string, skip, limit int) ([]*entity.SlackChannel, int64, error) {
	var err error
	var total int64
	var entities []*entity.SlackChannel

	err = r.db.
		Model(&entity.SlackChannel{}).
		Where("tenant_name = ?", tenant).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Model(&entity.SlackChannel{}).
		Offset(skip).
		Limit(limit).
		Where("tenant_name = ?", tenant).
		Find(&entities).Error

	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *slackChannelRepository) CreateSlackChannel(entity *entity.SlackChannel) error {
	err := r.db.Create(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *slackChannelRepository) UpdateSlackChannel(entityId uuid.UUID, organizationId string) error {
	return r.db.Model(&entity.SlackChannel{}).Where("id = ?", entityId).Update("organization_id", organizationId).Error
}
