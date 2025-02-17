package postgres_repository

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go/log"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type WebSessionRepository interface {
	Create(ctx context.Context, webWebSessionData postgres_entity.WebSession) (*postgres_entity.WebSession, error)
	FindAllActiveSessions(ctx context.Context, webWebSessionData postgres_entity.WebSession, sessionTimeoutInMins *int) ([]postgres_entity.WebSession, error)
	FindSession(ctx context.Context, webSessionData postgres_entity.WebSession, lookbackPeriodInMins *int) (*postgres_entity.WebSession, error)
	FindLastNotification(ctx context.Context, tenant, domain string) (*postgres_entity.WebSession, error)
	FindAllSessionsForSupportAnalysis(ctx context.Context) ([]postgres_entity.WebSession, error)
	UpdateLastActivity(ctx context.Context, sessionID, eventType string) (*postgres_entity.WebSession, error)
	UpdateSessionEnd(ctx context.Context, sessionID string, endTime time.Time) (*postgres_entity.WebSession, error)
	UpdateSessionWithDomain(ctx context.Context, sessionID, domain string) (*postgres_entity.WebSession, error)
	UpdateSessionWithOrganization(ctx context.Context, sessionID, organizationId string) error
	UpdateSessionWithSlackSentAt(ctx context.Context, sessionID string) error
	UpdateSessionPageViews(ctx context.Context, sessionID, tenant string, pageViews []string) (*postgres_entity.WebSession, error)
	UpdateSupportSignals(ctx context.Context, sessionID, tenant string, supportSignal int8) (*postgres_entity.WebSession, error)
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
	tracing.LogObjectAsJson(span, "webSessionData", webSessionData)

	var created postgres_entity.WebSession
	err := r.gormDb.Create(&webSessionData).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &created, nil
}

func (r *webSessionEventsRepository) FindAllActiveSessions(ctx context.Context, webSessionData postgres_entity.WebSession, sessionTimeoutInMins *int) ([]postgres_entity.WebSession, error) {
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

func (r *webSessionEventsRepository) FindAllSessionsForSupportAnalysis(ctx context.Context) ([]postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.FindAllSessionsForIntentAnalysis")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set on context")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var results []postgres_entity.WebSession
	err := r.gormDb.Model(&postgres_entity.WebSession{}).
		Where("tenant = ?", tenant).
		Where("support_signals = ?", postgres_entity.SupportNotAnalyzed).
		Where("is_active = ?", false).
		Where("organization_id IS NOT NULL AND organization_id != ''").
		Where("unique_page_views IS NOT NULL AND array_length(unique_page_views, 1) > 0").
		Order("created_at DESC").
		Find(&results).
		Error
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

	query := r.gormDb.Where(&webSessionData).Order("last_activity DESC")

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
	span.LogFields(log.String("tenant", tenant), log.String("domain", domain))

	var result postgres_entity.WebSession
	err := r.gormDb.
		Model(&postgres_entity.WebSession{}).
		Where("tenant = ? AND domain = ? AND sent_slack_notification IS NOT NULL", tenant, domain).
		Order("sent_slack_notification DESC").
		First(&result).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.LogFields(log.String("result", "No record found"))
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.String("result", "Record found: "+result.ID))
	return &result, nil
}

func (r *webSessionEventsRepository) UpdateLastActivity(ctx context.Context, sessionID, eventType string) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.UpdateLastActivity")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagEntity(span, sessionID)

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
	tracing.TagEntity(span, sessionID)

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
	tracing.TagEntity(span, sessionID)

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

func (r *webSessionEventsRepository) UpdateSessionWithOrganization(ctx context.Context, sessionID, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.UpdateSessionWithOrganization")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagEntity(span, sessionID)

	// Perform a single-column update on organization_id for the matching session
	result := r.gormDb.Model(&postgres_entity.WebSession{}).
		Where("id = ?", sessionID).
		Update("organization_id", &organizationId)

	if result.Error != nil {
		// Log the error with tracing and return it
		tracing.TraceErr(span, result.Error)
		return result.Error
	}

	return nil
}

func (r *webSessionEventsRepository) UpdateSessionPageViews(ctx context.Context, sessionID, tenant string, pageViews []string) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.UpdateSessionPageViews")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var updatedSession postgres_entity.WebSession
	err := r.gormDb.Model(&postgres_entity.WebSession{}).
		Where("id = ? AND tenant = ?", sessionID, tenant).
		Update("unique_page_views", pq.StringArray(pageViews)).
		First(&updatedSession).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedSession, nil
}

func (r *webSessionEventsRepository) UpdateSupportSignals(ctx context.Context, sessionID, tenant string, supportSignal int8) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.UpdateSessionPageViews")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagEntity(span, sessionID)

	if supportSignal != postgres_entity.SupportNeedDetected && supportSignal != postgres_entity.SupportNeedDetected {
		err := errors.New("invalid supportSignal value")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedSession postgres_entity.WebSession
	err := r.gormDb.Model(&postgres_entity.WebSession{}).
		Where("id = ? AND tenant = ?", sessionID, tenant).
		Update("support_signals", supportSignal).
		First(&updatedSession).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedSession, nil
}

func (r *webSessionEventsRepository) UpdateSessionWithSlackSentAt(ctx context.Context, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.UpdateSessionWithSlackSentAt")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagEntity(span, sessionID)

	// Perform a single-column update on sent_slack_notification for the matching session
	result := r.gormDb.Model(&postgres_entity.WebSession{}).
		Where("id = ?", sessionID).
		Update("sent_slack_notification", utils.Now())

	if result.Error != nil {
		// Log the error with tracing and return it
		tracing.TraceErr(span, result.Error)
		return result.Error
	}

	return nil
}
