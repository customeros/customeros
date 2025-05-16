package postgres_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type EmailTrackingRepository interface {
	Register(ctx context.Context, emailTracking postgres_entity.EmailTracking) (*postgres_entity.EmailTracking, error)
	Update(ctx context.Context, emailTracking postgres_entity.EmailTracking) (*postgres_entity.EmailTracking, error)
}

type emailTrackingRepository struct {
	gormDb *gorm.DB
}

func NewEmailTrackingRepository(gormDb *gorm.DB) EmailTrackingRepository {
	return &emailTrackingRepository{gormDb: gormDb}
}

func (e emailTrackingRepository) Register(ctx context.Context, emailTracking postgres_entity.EmailTracking) (*postgres_entity.EmailTracking, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EmailTrackingRepository.Register")
	defer spans.Finish()
	spans.LogObjectAsJson("emailTracking", emailTracking)

	err := e.gormDb.Create(&emailTracking).Error

	if err != nil {
		spans.TraceError(err)
		return nil, errors.Wrap(err, "failed to store email lookup")
	}

	return &emailTracking, nil
}

func (e emailTrackingRepository) Update(ctx context.Context, emailTracking postgres_entity.EmailTracking) (*postgres_entity.EmailTracking, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EmailTrackingRepository.Update")
	defer spans.Finish()

	// Fetch the existing record
	var existingTracking postgres_entity.EmailTracking
	if err := e.gormDb.First(&existingTracking, "id = ?", emailTracking.ID).Error; err != nil {
		spans.TraceError(err)
		return nil, errors.Wrap(err, "failed to find email tracking record")
	}

	// Update only specific fields
	updates := map[string]interface{}{
		"updated_at": utils.Now(),
		"ip":         emailTracking.IP,
	}

	err := e.gormDb.Model(&existingTracking).Updates(updates).Error
	if err != nil {
		spans.TraceError(err)
		return nil, errors.Wrap(err, "failed to update email tracking")
	}

	// Refresh the struct with updated data
	if err := e.gormDb.First(&existingTracking, existingTracking.ID).Error; err != nil {
		spans.TraceError(err)
		return nil, errors.Wrap(err, "failed to refresh email tracking data")
	}

	return &existingTracking, nil
}
