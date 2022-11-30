package repository

import (
	"database/sql"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/entity"
	"gorm.io/gorm"
)

type PageViewRepository interface {
	GetPageViewsForSync(bucketSize int) (entity.PageViews, error)
	MarkSynced(pv entity.PageView, contactId string) error
}

type pageViewRepository struct {
	db *gorm.DB
}

func NewPageViewRepository(gormDb *gorm.DB) PageViewRepository {
	return &pageViewRepository{
		db: gormDb,
	}
}

func (r *pageViewRepository) GetPageViewsForSync(bucketSize int) (entity.PageViews, error) {
	var pageViews entity.PageViews

	err := r.db.Limit(bucketSize).
		Where("visitor_id is not null and visitor_id <> ? and synced_to_customer_os = false", "").
		Find(&pageViews).Error
	if err != nil {
		return nil, err
	}

	return pageViews, nil
}

func (r *pageViewRepository) MarkSynced(pv entity.PageView, contactId string) error {
	return r.db.Model(&pv).
		Updates(entity.PageView{
			SyncedToCustomerOs: true,
			ContactID:          sql.NullString{String: contactId, Valid: true},
		}).
		Error
}
