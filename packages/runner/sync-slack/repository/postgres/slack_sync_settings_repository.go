package postgresrepo

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
	"time"
)

type SlackSyncSettingsRepository interface {
	AutoMigrate() error
	GetChannelsToSync(ctx context.Context) ([]entity.SlackSyncSettings, error)
	SaveSyncRun(ctx context.Context, slackStngs entity.SlackSyncSettings, at time.Time) error
}

type slackSyncRepository struct {
	db *gorm.DB
}

func NewSlackSyncSettingsRepository(db *gorm.DB) SlackSyncSettingsRepository {
	return &slackSyncRepository{
		db: db,
	}
}

func (r *slackSyncRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&entity.SlackSyncSettings{})
}

func (r *slackSyncRepository) GetChannelsToSync(ctx context.Context) ([]entity.SlackSyncSettings, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackSyncSettingsRepository.GetChannelsToSync")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(span)

	var settings []entity.SlackSyncSettings

	err := r.db.Where("organization_id is not null AND organization_id != ? AND slack_access = ? AND enabled = ?",
		"", true, true).Find(&settings).Error

	if err != nil {
		return nil, err
	}

	return settings, nil

}

func (r *slackSyncRepository) SaveSyncRun(ctx context.Context, slackStngs entity.SlackSyncSettings, at time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackSyncSettingsRepository.SaveSyncRun")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(span)
	span.LogFields(log.Object("slackSettings", slackStngs), log.Object("at", at))

	slackSync := entity.SlackSyncSettings{
		Tenant:    slackStngs.Tenant,
		ChannelId: slackStngs.ChannelId,
	}
	r.db.FirstOrCreate(&slackSync, slackSync)

	slackSync.LastSyncAt = &at
	if slackStngs.TeamId != "" {
		slackSync.TeamId = slackStngs.TeamId
	}
	if slackStngs.ChannelName != "" {
		slackSync.ChannelName = slackStngs.ChannelName
	}

	return r.db.Model(&slackSync).
		Where(&entity.SlackSyncSettings{Tenant: slackStngs.Tenant, ChannelId: slackStngs.ChannelId}).
		Save(&slackSync).Error
}
