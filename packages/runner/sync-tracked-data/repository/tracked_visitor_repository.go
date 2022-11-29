package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/entity"
	"gorm.io/gorm"
)

type TrackedVisitorRepository interface {
	GetVisitorsWithSyncedFalse(context.Context) ([]entity.Visitor, error)
}

type trackedVisitorRepository struct {
	db *gorm.DB
}

func NewTrackedVisitorRepository(gormDb *gorm.DB) TrackedVisitorRepository {
	return &trackedVisitorRepository{
		db: gormDb,
	}
}

func (r *trackedVisitorRepository) GetVisitorsWithSyncedFalse(ctx context.Context) ([]entity.Visitor, error) {
	var visitors entity.Visitors

	err := r.db.Where("visitor_id is not null and synced_to_customer_os = false").Find(&visitors).Error

	if err != nil {
		return nil, err
	}

	return visitors, nil
}
