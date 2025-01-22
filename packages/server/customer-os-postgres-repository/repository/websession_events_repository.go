package postgres_repository

import (
	"context"
	"errors"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type WebSessionRepository interface {
	Create(ctx context.Context, webWebSessionData postgres_entity.WebSession) (*postgres_entity.WebSession, error)
	FindAllSessions(ctx context.Context, webWebSessionData postgres_entity.WebSession, sessionTimeoutInMins *int) ([]postgres_entity.WebSession, error)
	FindSession(ctx context.Context, webSessionData postgres_entity.WebSession, lookbackPeriodInMins *int) (*postgres_entity.WebSession, error)
	FindLastNotification(ctx context.Context, tenant, domain string) (*postgres_entity.WebSession, error)
	UpdateLastActivity(ctx context.Context, sessionID, eventType string) (*postgres_entity.WebSession, error)
	UpdateSessionEnd(ctx context.Context, sessionID string, endTime time.Time) (*postgres_entity.WebSession, error)
	UpdateSessionWithDomain(ctx context.Context, sessionID, domain string) (*postgres_entity.WebSession, error)
}

type webSessionEventsRepository struct {
	gormDb *gorm.DB
}

func NewWebSessionRepository(gormDb *gorm.DB) WebSessionRepository {
	return &webSessionEventsRepository{gormDb: gormDb}
}

func (r *webSessionEventsRepository) Create(ctx context.Context, webSessionData postgres_entity.WebSession) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var created postgres_entity.WebSession
	err := r.gormDb.Create(&webSessionData).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &created, nil
}

func (r *webSessionEventsRepository) FindAllSessions(ctx context.Context, webSessionData postgres_entity.WebSession, sessionTimeoutInMins *int) ([]postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// Start with a base query
	query := r.gormDb.Model(&postgres_entity.WebSession{})

	// Add conditions explicitly
	if webSessionData.Tenant != "" {
		query = query.Where("tenant = ?", webSessionData.Tenant)
	}
	if webSessionData.Domain != nil {
		query = query.Where("domain = ?", *webSessionData.Domain)
	}
	if webSessionData.VisitorID != "" {
		query = query.Where("visitor_id = ?", webSessionData.VisitorID)
	}

	// Explicitly add the is_active condition
	query = query.Where("is_active = ?", webSessionData.IsActive)

	// Add lookback period if provided
	if sessionTimeoutInMins != nil {
		lookbackDate := time.Now().Add(-time.Duration(*sessionTimeoutInMins) * time.Minute)
		query = query.Where("last_activity < ?", lookbackDate)
	}

	// Order and execute
	var results []postgres_entity.WebSession
	err := query.Order("created_at DESC").Find(&results).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return results, nil
}

func (r *webSessionEventsRepository) FindSession(ctx context.Context, webSessionData postgres_entity.WebSession, lookbackPeriodInMins *int) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	query := r.gormDb.Where(&webSessionData).Order("created_at DESC")

	// Add lookback period if provided
	if lookbackPeriodInMins != nil {
		lookbackDate := time.Now().Add(-time.Duration(*lookbackPeriodInMins) * time.Minute)
		query = query.Where("last_activity > ?", lookbackDate)
	}

	var result postgres_entity.WebSession
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

func (r *webSessionEventsRepository) FindLastNotification(ctx context.Context, tenant, domain string) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.FindLastNotification")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var result postgres_entity.WebSession
	err := r.gormDb.
		Model(&postgres_entity.WebSession{}).
		Where("tenant = ? AND domain = ?", tenant, domain).
		Order("sent_slack_notification DESC").
		First(&result).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &result, nil
}

func (r *webSessionEventsRepository) UpdateLastActivity(ctx context.Context, sessionID, eventType string) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.UpdateLastActivity")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var updatedSession postgres_entity.WebSession
	err := r.gormDb.
		Model(&postgres_entity.WebSession{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"last_activity":   utils.Now(),
			"last_event_type": eventType,
		}).
		First(&updatedSession, "id = ?", sessionID).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedSession, nil
}

func (r *webSessionEventsRepository) UpdateSessionEnd(ctx context.Context, sessionID string, endTime time.Time) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.UpdateSessionEnd")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var updatedSession postgres_entity.WebSession
	err := r.gormDb.
		Model(&postgres_entity.WebSession{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"is_active":       false,
			"end_time":        endTime,
			"published_event": true,
		}).
		First(&updatedSession, "id = ?", sessionID).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedSession, nil
}

func (r *webSessionEventsRepository) UpdateSessionWithDomain(ctx context.Context, sessionID, domain string) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.UpdateSessionWithDomain")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var updatedSession postgres_entity.WebSession
	err := r.gormDb.Model(&postgres_entity.WebSession{}).
		Where("id = ?", sessionID).
		Update("domain", &domain).
		First(&updatedSession).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedSession, nil
}
