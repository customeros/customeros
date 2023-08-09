package postgresrepo

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/repository/helper"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
	"time"
)

type SlackSyncRepository interface {
	AutoMigrate() error
	FindForTenantAndChannelId(ctx context.Context, tenant, channelId string) helper.QueryResult
	SaveSyncRun(ctx context.Context, tenant, channelId string, at time.Time) error
}

type slackSyncRepository struct {
	db *gorm.DB
}

func NewSlackSyncRepository(db *gorm.DB) SlackSyncRepository {
	return &slackSyncRepository{
		db: db,
	}
}

func (r *slackSyncRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&entity.SlackSync{})
}

func (r *slackSyncRepository) FindForTenantAndChannelId(ctx context.Context, tenant, channelId string) helper.QueryResult {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackSyncRepository.FindForTenantAndChannelId")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(span)
	span.LogFields(log.String("tenant", tenant), log.String("channelId", channelId))

	var slackSync entity.SlackSync

	err := r.db.
		Where("tenant = ? and channel_id = ?", tenant, channelId).
		First(&slackSync).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return helper.QueryResult{Error: err}
	}
	if err == gorm.ErrRecordNotFound {
		return helper.QueryResult{Result: nil}
	}

	return helper.QueryResult{Result: slackSync}
}

func (r *slackSyncRepository) SaveSyncRun(ctx context.Context, tenant, channelId string, at time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackSyncRepository.SaveSyncRun")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(span)
	span.LogFields(log.String("tenant", tenant), log.String("channelId", channelId), log.Object("at", at))

	slackSync := entity.SlackSync{
		Tenant:    tenant,
		ChannelId: channelId,
	}
	r.db.FirstOrCreate(&slackSync, slackSync)
	slackSync.LastSyncAt = at
	return r.db.Model(&slackSync).
		Where(&entity.SlackSync{Tenant: tenant, ChannelId: channelId}).
		Save(&slackSync).Error
}
