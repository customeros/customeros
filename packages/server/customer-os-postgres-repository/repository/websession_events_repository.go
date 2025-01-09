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

type WebSessionRepository interface {
	Create(ctx context.Context, webWebSessionData entity.WebSession) (*entity.WebSession, error)
	FindAllTimedOutSessions(ctx context.Context, webWebSessionData entity.WebSession, sessionTimeoutInMins *int) ([]entity.WebSession, error)
	FindSession(ctx context.Context, webSessionData entity.WebSession, lookbackPeriodInMins *int) (*entity.WebSession, error)
	Update(ctx context.Context, webSessionData entity.WebSession) (*entity.WebSession, error)
}

type webSessionEventsRepository struct {
	gormDb *gorm.DB
}

func NewWebSessionRepository(gormDb *gorm.DB) WebSessionRepository {
	return &webSessionEventsRepository{gormDb: gormDb}
}

func (r *webSessionEventsRepository) Create(ctx context.Context, webSessionData entity.WebSession) (*entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var created entity.WebSession
	err := r.gormDb.Create(&webSessionData).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &created, nil
}

func (r *webSessionEventsRepository) FindAllTimedOutSessions(ctx context.Context, webSessionData entity.WebSession, sessionTimeoutInMins *int) ([]entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	query := r.gormDb.Where(&webSessionData).Order("created_at DESC")

	// Add lookback period if provided
	if sessionTimeoutInMins != nil {
		lookbackDate := time.Now().Add(-time.Duration(*sessionTimeoutInMins) * time.Minute)
		query = query.Where("last_activity < ?", lookbackDate)
	}

	var results []entity.WebSession
	err := query.Find(&results).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return results, nil
}

func (r *webSessionEventsRepository) FindSession(ctx context.Context, webSessionData entity.WebSession, lookbackPeriodInMins *int) (*entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	query := r.gormDb.Where(&webSessionData).Order("created_at DESC")

	// Add lookback period if provided
	if lookbackPeriodInMins != nil {
		lookbackDate := time.Now().Add(-time.Duration(*lookbackPeriodInMins) * time.Minute)
		query = query.Where("last_activity > ?", lookbackDate)
	}

	var result entity.WebSession
	err := query.First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &result, nil
}

func (r *webSessionEventsRepository) Update(ctx context.Context, webSessionData entity.WebSession) (*entity.WebSession, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "WebSessionRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if webSessionData.ID == "" {
		err := errors.New("ID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedSession entity.WebSession
	err := r.gormDb.
		Model(&entity.WebSession{}).
		Where("id = ?", webSessionData.ID).
		Updates(&webSessionData).
		First(&updatedSession).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedSession, nil
}
