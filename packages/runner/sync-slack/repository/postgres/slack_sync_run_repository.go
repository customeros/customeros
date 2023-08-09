package postgresrepo

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/entity"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

type SlackSyncRunRepository interface {
	AutoMigrate() error
	Save(ctx context.Context, entity entity.SlackSyncRunStatus)
}

type slackSyncRunRepository struct {
	db *gorm.DB
}

func NewSlackSyncRunRepository(gormDb *gorm.DB) SlackSyncRunRepository {
	return &slackSyncRunRepository{
		db: gormDb,
	}
}

func (r *slackSyncRunRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&entity.SlackSyncRunStatus{})
}

func (r *slackSyncRunRepository) Save(ctx context.Context, entity entity.SlackSyncRunStatus) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackSyncRunRepository.Save")
	defer span.Finish()

	r.db.Create(&entity)
}
