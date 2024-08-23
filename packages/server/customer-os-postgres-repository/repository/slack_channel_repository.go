package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type slackChannelRepository struct {
	db *gorm.DB
}

type SlackChannelRepository interface {
	GetSlackChannel(ctx context.Context, tenant, channelId string) (*entity.SlackChannel, error)
	GetSlackChannels(ctx context.Context, tenant string) ([]*entity.SlackChannel, error)
	GetPaginatedSlackChannels(ctx context.Context, tenant string, skip, limit int) ([]*entity.SlackChannel, int64, error)

	CreateSlackChannel(ctx context.Context, entity *entity.SlackChannel) error
	UpdateSlackChannelOrganization(ctx context.Context, entityId uuid.UUID, organizationId string) error
	UpdateSlackChannelName(ctx context.Context, entityId uuid.UUID, channelName string) error
}

func NewSlackChannelRepository(db *gorm.DB) SlackChannelRepository {
	return &slackChannelRepository{db: db}
}

func (r *slackChannelRepository) GetSlackChannel(ctx context.Context, tenant, channelId string) (*entity.SlackChannel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackChannelRepository.GetSlackChannel")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

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

func (r *slackChannelRepository) GetSlackChannels(ctx context.Context, tenant string) ([]*entity.SlackChannel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackChannelRepository.GetSlackChannels")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	var entities []*entity.SlackChannel
	err := r.db.
		Where("tenant_name = ?", tenant).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *slackChannelRepository) GetPaginatedSlackChannels(ctx context.Context, tenant string, skip, limit int) ([]*entity.SlackChannel, int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackChannelRepository.GetPaginatedSlackChannels")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

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

func (r *slackChannelRepository) CreateSlackChannel(ctx context.Context, entity *entity.SlackChannel) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackChannelRepository.CreateSlackChannel")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.db.Create(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *slackChannelRepository) UpdateSlackChannelOrganization(ctx context.Context, entityId uuid.UUID, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackChannelRepository.UpdateSlackChannelOrganization")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	return r.db.Model(&entity.SlackChannel{}).Where("id = ?", entityId).Update("organization_id", organizationId).Error
}

func (r *slackChannelRepository) UpdateSlackChannelName(ctx context.Context, entityId uuid.UUID, channelName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackChannelRepository.UpdateSlackChannelName")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	return r.db.Model(&entity.SlackChannel{}).Where("id = ?", entityId).Update("channel_name", channelName).Error
}
