package postgres_repository

import (
	"context"
	"errors"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type WebSessionRepository interface {
	Create(ctx context.Context, webWebSessionData postgres_entity.WebSession) (*postgres_entity.WebSession, error)
	FindAllActiveSessions(ctx context.Context, webWebSessionData postgres_entity.WebSession, sessionTimeoutInMins *int) ([]postgres_entity.WebSession, error)
	FindSession(ctx context.Context, webSessionData postgres_entity.WebSession, lookbackPeriodInMins *int) (*postgres_entity.WebSession, error)
	FindLastNotification(ctx context.Context, tenant, domain string) (*postgres_entity.WebSession, error)
	FindLatestSessionWithDomainByIP(ctx context.Context, ipAddress string) (*postgres_entity.WebSession, error)
	UpdateLastActivity(ctx context.Context, sessionID, eventType string) (*postgres_entity.WebSession, error)
	SetSessionEnd(ctx context.Context, sessionID string, endTime time.Time) (*postgres_entity.WebSession, error)
	SetOrganizationId(ctx context.Context, sessionID, organizationId string) error
	SetSlackSentAt(ctx context.Context, sessionID string) error
	SetSessionPageViews(ctx context.Context, sessionID, tenant string, pageViews []string) (*postgres_entity.WebSession, error)
	SetVisitorIdentity(ctx context.Context, sessionID string, domain, email *string, emailType *enum.EmailType) error
}

type webSessionRepository struct {
	gormDb *gorm.DB
}

func NewWebSessionRepository(gormDb *gorm.DB) WebSessionRepository {
	return &webSessionRepository{gormDb: gormDb}
}

func (r *webSessionRepository) Create(ctx context.Context, webSessionData postgres_entity.WebSession) (*postgres_entity.WebSession, error) {
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

func (r *webSessionRepository) FindAllActiveSessions(ctx context.Context, webSessionData postgres_entity.WebSession, sessionTimeoutInMins *int) ([]postgres_entity.WebSession, error) {
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

func (r *webSessionRepository) FindSession(ctx context.Context, webSessionData postgres_entity.WebSession, lookbackPeriodInMins *int) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.Find")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

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
			span.LogFields(log.Bool("result.found", false))
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.String("result.sessionId", result.ID))
	return &result, nil
}

func (r *webSessionRepository) FindLastNotification(ctx context.Context, tenant, domain string) (*postgres_entity.WebSession, error) {
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

func (r *webSessionRepository) FindLatestSessionWithDomainByIP(ctx context.Context, ipAddress string) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.FindLatestSessionWithDomainByIP")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	// Build the query with IP and ensure domain is not null
	query := r.gormDb.Where("ip = ?", ipAddress).
		Where("domain IS NOT NULL").
		Where("domain != ''"). // Ensure domain is not empty string
		Order("last_activity DESC")

	var result postgres_entity.WebSession
	err := query.First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.LogFields(log.Bool("result.found", false))
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.String("result.sessionId", result.ID))
	span.LogFields(log.String("result.domain", *result.Domain))
	return &result, nil
}

func (r *webSessionRepository) UpdateLastActivity(ctx context.Context, sessionID, eventType string) (*postgres_entity.WebSession, error) {
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

func (r *webSessionRepository) SetSessionEnd(ctx context.Context, sessionID string, endTime time.Time) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.SetSessionEnd")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
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

func (r *webSessionRepository) SetOrganizationId(ctx context.Context, sessionID, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.SetOrganizationId")
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

func (r *webSessionRepository) SetSessionPageViews(ctx context.Context, sessionID, tenant string, pageViews []string) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.SetSessionPageViews")
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

func (r *webSessionRepository) SetSlackSentAt(ctx context.Context, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.SetSlackSentAt")
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

func (r *webSessionRepository) SetVisitorIdentity(ctx context.Context, sessionID string, domain, email *string, emailType *enum.EmailType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebSessionRepository.SetVisitorIdentity")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// Start building the query
	query := r.gormDb.Model(&postgres_entity.WebSession{}).
		Where("id = ?", sessionID)

	// Build a map of fields to update, only including non-nil values
	updates := map[string]interface{}{}

	if domain != nil && *domain != "" {
		updates["domain"] = *domain
	}

	if email != nil && *email != "" {
		updates["email"] = *email
	}

	if emailType != nil && *emailType != "" {
		updates["email_type"] = *emailType
	}

	// If we have fields to update, execute the update
	if len(updates) > 0 {
		err := query.Updates(updates).Error
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
