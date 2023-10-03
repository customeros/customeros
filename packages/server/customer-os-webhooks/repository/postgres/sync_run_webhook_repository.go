package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository/postgres/entity"
	"gorm.io/gorm"
)

type SyncRunWebhookRepository interface {
	Save(ctx context.Context, entity postgresentity.SyncRunWebhook)
}

type syncRunWebhookRepository struct {
	db *gorm.DB
}

func NewSyncRunWebhookRepository(gormDb *gorm.DB) SyncRunWebhookRepository {
	return &syncRunWebhookRepository{
		db: gormDb,
	}
}

func (r *syncRunWebhookRepository) Save(ctx context.Context, entity postgresentity.SyncRunWebhook) {
	r.db.Create(&entity)
}
