package postgres_repository

import (
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type slackChannelRepository struct {
	db *gorm.DB
}

type SlackChannelRepository interface {
	GetSlackChannel(ctx context.Context, tenant, channelId string) (*postgres_entity.SlackChannel, error)
	GetSlackChannels(ctx context.Context, tenant string) ([]*postgres_entity.SlackChannel, error)
	GetPaginatedSlackChannels(ctx context.Context, tenant string, skip, limit int) ([]*postgres_entity.SlackChannel, int64, error)

	CreateSlackChannel(ctx context.Context, postgres_entity *postgres_entity.SlackChannel) error
	UpdateSlackChannelOrganization(ctx context.Context, postgres_entityId uuid.UUID, organizationId string) error
	UpdateSlackChannelName(ctx context.Context, postgres_entityId uuid.UUID, channelName string) error
}

func NewSlackChannelRepository(db *gorm.DB) SlackChannelRepository {
	return &slackChannelRepository{db: db}
}

func (r *slackChannelRepository) GetSlackChannel(ctx context.Context, tenant, channelId string) (*postgres_entity.SlackChannel, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "SlackChannelRepository.GetSlackChannel")
	defer spans.Finish()
	spans.LogKV("channelId", channelId)

	var entities []postgres_entity.SlackChannel
	err := r.db.
		Where("tenant_name = ?", tenant).
		Where("channel_id = ?", channelId).
		Find(&entities).Error

	if err != nil {
		spans.TraceError(err)
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

func (r *slackChannelRepository) GetSlackChannels(ctx context.Context, tenant string) ([]*postgres_entity.SlackChannel, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "SlackChannelRepository.GetSlackChannels")
	defer spans.Finish()
	spans.LogKV("tenant", tenant)

	var entities []*postgres_entity.SlackChannel
	err := r.db.
		Where("tenant_name = ?", tenant).
		Find(&entities).Error

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return entities, nil
}

func (r *slackChannelRepository) GetPaginatedSlackChannels(ctx context.Context, tenant string, skip, limit int) ([]*postgres_entity.SlackChannel, int64, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "SlackChannelRepository.GetPaginatedSlackChannels")
	defer spans.Finish()
	spans.LogKV("tenant", tenant, "skip", skip, "limit", limit)

	var err error
	var total int64
	var entities []*postgres_entity.SlackChannel

	err = r.db.
		Model(&postgres_entity.SlackChannel{}).
		Where("tenant_name = ?", tenant).
		Count(&total).Error
	if err != nil {
		spans.TraceError(err)
		return nil, 0, err
	}

	err = r.db.
		Model(&postgres_entity.SlackChannel{}).
		Offset(skip).
		Limit(limit).
		Where("tenant_name = ?", tenant).
		Find(&entities).Error

	if err != nil {
		spans.TraceError(err)
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *slackChannelRepository) CreateSlackChannel(ctx context.Context, postgres_entity *postgres_entity.SlackChannel) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "SlackChannelRepository.CreateSlackChannel")
	defer spans.Finish()

	err := r.db.Create(postgres_entity).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *slackChannelRepository) UpdateSlackChannelOrganization(ctx context.Context, postgres_entityId uuid.UUID, organizationId string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "SlackChannelRepository.UpdateSlackChannelOrganization")
	defer spans.Finish()
	spans.LogKV("postgres_entityId", postgres_entityId, "organizationId", organizationId)

	return r.db.Model(&postgres_entity.SlackChannel{}).Where("id = ?", postgres_entityId).Update("organization_id", organizationId).Error
}

func (r *slackChannelRepository) UpdateSlackChannelName(ctx context.Context, postgres_entityId uuid.UUID, channelName string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "SlackChannelRepository.UpdateSlackChannelName")
	defer spans.Finish()
	spans.LogKV("postgres_entityId", postgres_entityId, "channelName", channelName)

	return r.db.Model(&postgres_entity.SlackChannel{}).Where("id = ?", postgres_entityId).Update("channel_name", channelName).Error
}
