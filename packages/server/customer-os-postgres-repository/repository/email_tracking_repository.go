package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type EmailTrackingRepository interface {
	Register(ctx context.Context, emailTracking entity.EmailTracking) (*entity.EmailTracking, error)
	Update(ctx context.Context, emailTracking entity.EmailTracking) (*entity.EmailTracking, error)
}

type emailTrackingRepository struct {
	gormDb *gorm.DB
}

func NewEmailTrackingRepository(gormDb *gorm.DB) EmailTrackingRepository {
	return &emailTrackingRepository{gormDb: gormDb}
}

func (e emailTrackingRepository) Register(ctx context.Context, emailTracking entity.EmailTracking) (*entity.EmailTracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailTrackingRepository.Register")
	defer span.Finish()
	emailTracking.ID = utils.GenerateRandomString(64)

	err := e.gormDb.Create(&emailTracking).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "failed to store email lookup")
	}

	return &emailTracking, nil
}

func (e emailTrackingRepository) Update(ctx context.Context, emailTracking entity.EmailTracking) (*entity.EmailTracking, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailTrackingRepository.Update")
	defer span.Finish()

	// Fetch the existing record
	var existingTracking entity.EmailTracking
	if err := e.gormDb.First(&existingTracking, "id = ?", emailTracking.ID).Error; err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "failed to find email tracking record")
	}

	// Update only specific fields
	updates := map[string]interface{}{
		"updated_at": time.Now(),
		"ip":         emailTracking.IP,
	}

	err := e.gormDb.Model(&existingTracking).Updates(updates).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "failed to update email tracking")
	}

	// Refresh the struct with updated data
	if err := e.gormDb.First(&existingTracking, existingTracking.ID).Error; err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "failed to refresh email tracking data")
	}

	return &existingTracking, nil
}
