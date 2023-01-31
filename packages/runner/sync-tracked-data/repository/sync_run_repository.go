package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/entity"
	"gorm.io/gorm"
)

type SyncRunRepository interface {
	Save(entity entity.SyncRun)
}

type syncRunRepository struct {
	db *gorm.DB
}

func NewSyncRunRepository(gormDb *gorm.DB) SyncRunRepository {
	return &syncRunRepository{
		db: gormDb,
	}
}

func (r *syncRunRepository) Save(entity entity.SyncRun) {
	r.db.Create(&entity)
}
