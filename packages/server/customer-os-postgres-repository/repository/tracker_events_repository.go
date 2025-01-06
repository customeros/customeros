package repository

import (
	"context"
	"errors"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type TrackerEventsRepository interface {
	Create(ctx context.Context, trackerData entity.TrackerEvents) (*entity.TrackerEvents, error)
	FindAll(ctx context.Context, trackerData entity.TrackerEvents, cacheLookbackInDays *int) (*[]entity.TrackerEvents, error)
	Update(ctx context.Context, trackerData entity.TrackerEvents) (*entity.TrackerEvents, error)
	UpdateWhereNoCompanyID(ctx context.Context, trackerData entity.TrackerEvents) (*entity.TrackerEvents, error)
}

type trackerEventsRepository struct {
	gormDb *gorm.DB
}

func NewTrackerEventsRepository(gormDb *gorm.DB) TrackerEventsRepository {
	return &trackerEventsRepository{gormDb: gormDb}
}

func (r *trackerEventsRepository) Create(ctx context.Context, trackerData entity.TrackerEvents) (*entity.TrackerEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackerEventsRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var created entity.TrackerEvents
	err := r.gormDb.Create(&trackerData).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &created, nil
}

func (r *trackerEventsRepository) FindAll(ctx context.Context, trackerData entity.TrackerEvents, cacheLookbackInDays *int) (*[]entity.TrackerEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackerEventsRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	query := r.gormDb.Where(&trackerData).Order("created_at DESC")

	// Add lookback period if provided
	if cacheLookbackInDays != nil {
		lookbackDate := time.Now().AddDate(0, 0, -*cacheLookbackInDays)
		query = query.Where("created_at > ?", lookbackDate)
	}

	var results []entity.TrackerEvents
	err := query.Find(&results).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &results, nil
}

func (f *trackerEventsRepository) Update(ctx context.Context, trackerData entity.TrackerEvents) (*entity.TrackerEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackerEventsRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if trackerData.IP == "" && trackerData.VisitorId == "" {
		err := errors.New("visitor ID or IP address is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedRecord entity.TrackerEvents
	err := f.gormDb.Model(&trackerData).Updates(&trackerData).First(&updatedRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedRecord, nil
}

func (f *trackerEventsRepository) UpdateWhereNoCompanyID(ctx context.Context, trackerData entity.TrackerEvents) (*entity.TrackerEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackerEventsRepository.UpdateWhereNoCompanyID")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if trackerData.IP == "" && trackerData.VisitorId == "" {
		err := errors.New("visitor ID or IP address is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedRecord entity.TrackerEvents
	err := f.gormDb.
		Model(&trackerData).
		Where("(domain IS NULL OR domain = '')").
		Updates(&trackerData).
		First(&updatedRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedRecord, nil
}
