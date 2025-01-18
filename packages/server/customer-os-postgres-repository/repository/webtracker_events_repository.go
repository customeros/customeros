package repository

import (
	"context"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type WebTrackerEventsRepository interface {
	Create(ctx context.Context, webWebTrackerData entity.WebTrackerEvents) (*entity.WebTrackerEvents, error)
	FindAll(ctx context.Context, webWebTrackerData entity.WebTrackerEvents, cacheLookbackInDays *int) ([]entity.WebTrackerEvents, error)
}

type webTrackerEventsRepository struct {
	gormDb *gorm.DB
}

func NewWebTrackerEventsRepository(gormDb *gorm.DB) WebTrackerEventsRepository {
	return &webTrackerEventsRepository{gormDb: gormDb}
}

func (r *webTrackerEventsRepository) Create(ctx context.Context, webTrackerData entity.WebTrackerEvents) (*entity.WebTrackerEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebTrackerEventsRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var created entity.WebTrackerEvents
	err := r.gormDb.Create(&webTrackerData).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &created, nil
}

func (r *webTrackerEventsRepository) FindAll(ctx context.Context, webTrackerData entity.WebTrackerEvents, cacheLookbackInDays *int) ([]entity.WebTrackerEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebTrackerEventsRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	query := r.gormDb.Where(&webTrackerData).Order("created_at DESC")

	// Add lookback period if provided
	if cacheLookbackInDays != nil {
		lookbackDate := time.Now().AddDate(0, 0, -*cacheLookbackInDays)
		query = query.Where("created_at > ?", lookbackDate)
	}

	var results []entity.WebTrackerEvents
	err := query.Find(&results).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return results, nil
}
