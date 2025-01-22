package postgres_repository

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
	"time"
)

type GlobalOrganizationRepository interface {
	GetById(ctx context.Context, id uint64) (*postgres_entity.GlobalOrganization, error)
	GetByPrimaryDomain(ctx context.Context, domain string) (*postgres_entity.GlobalOrganization, error)
	GetByPrimaryDomains(ctx context.Context, domains []string) ([]*postgres_entity.GlobalOrganization, error)
	Create(ctx context.Context, organization *postgres_entity.GlobalOrganization) (*postgres_entity.GlobalOrganization, error)
	Update(ctx context.Context, organization *postgres_entity.GlobalOrganization) (*postgres_entity.GlobalOrganization, error)
	Search(ctx context.Context, searchTerm string, limit int) ([]*postgres_entity.GlobalOrganization, error)
	GetOrganizationsToEnrichIndustry(ctx context.Context, hoursFromPreviousAttempt, maxAttempts, limit int) ([]*postgres_entity.GlobalOrganization, error)
	GetOrganizationsToEnrichDescription(ctx context.Context, hoursFromPreviousAttempt, maxAttempts, limit int) ([]*postgres_entity.GlobalOrganization, error)
	GetOrganizationsToEnrichName(ctx context.Context, hoursFromPreviousAttempt, maxAttempts, limit int) ([]*postgres_entity.GlobalOrganization, error)
	MarkIndustryEnrichRequested(ctx context.Context, id uint64) error
	MarkDescriptionEnrichRequested(ctx context.Context, id uint64) error
	MarkNameEnrichRequested(ctx context.Context, id uint64) error
	SetIndustry(ctx context.Context, id uint64, industryNaicsCode, industryNaicsName string) error
	SetDescription(ctx context.Context, id uint64, description string) error
	SetName(ctx context.Context, id uint64, name string) error
	GetGlobalOrganizationsToSyncIntoTenantOrganizations(ctx context.Context, daysFromPreviousSync, limit int) ([]*postgres_entity.GlobalOrganization, error)
	MarkGlobalOrganizationSyncedToNeo(ctx context.Context, id uint64) error
}

type globalOrganizationRepository struct {
	db *gorm.DB
}

func NewGlobalOrganizationRepository(gormDb *gorm.DB) GlobalOrganizationRepository {
	return &globalOrganizationRepository{db: gormDb}
}

func (r *globalOrganizationRepository) GetById(ctx context.Context, id uint64) (*postgres_entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Uint64("id", id))

	organization := &postgres_entity.GlobalOrganization{}
	result := r.db.WithContext(ctx).Where("id = ?", id).First(organization)
	if result.Error != nil {
		span.LogFields(tracingLog.Bool("result.found", false))
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogFields(tracingLog.Bool("result.found", true))
	return organization, nil
}

func (r *globalOrganizationRepository) GetByPrimaryDomains(ctx context.Context, domains []string) ([]*postgres_entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.GetByPrimaryDomains")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Object("domains", domains))

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).Where("primary_domain IN ?", domains).Find(&organizations)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogFields(tracingLog.Int("result.count", len(organizations)))
	return organizations, nil
}

func (r *globalOrganizationRepository) GetByPrimaryDomain(ctx context.Context, domain string) (*postgres_entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.GetByPrimaryDomain")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("domain", domain))

	organization := &postgres_entity.GlobalOrganization{}
	result := r.db.WithContext(ctx).Where("primary_domain = ?", domain).First(organization)
	if result.Error != nil {
		span.LogFields(tracingLog.Bool("found", false))
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogFields(tracingLog.Bool("found", true))
	return organization, nil
}

func (r *globalOrganizationRepository) Create(ctx context.Context, organization *postgres_entity.GlobalOrganization) (*postgres_entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "organization", organization)

	result := r.db.WithContext(ctx).Create(&organization)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	return organization, nil
}

func (r *globalOrganizationRepository) Update(ctx context.Context, organization *postgres_entity.GlobalOrganization) (*postgres_entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "organization", organization)

	organization.UpdatedAt = utils.Now()
	result := r.db.WithContext(ctx).Save(organization)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	return organization, nil
}

func (r *globalOrganizationRepository) Search(ctx context.Context, searchTerm string, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.SearchOrganizations")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("searchTerm", searchTerm))

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Select("id, name, primary_domain, website, logo_url, icon_url, other_domains").
		Where("name ILIKE ? OR primary_domain ILIKE ? OR other_domains ILIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogFields(tracingLog.Int("found", len(organizations)))
	return organizations, nil
}

func (r *globalOrganizationRepository) GetOrganizationsToEnrichIndustry(ctx context.Context, hoursFromPreviousAttempt, maxAttempts, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.GetOrganizationsToEnrichIndustry")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Int("hoursFromPreviousAttempt", hoursFromPreviousAttempt), tracingLog.Int("limit", limit), tracingLog.Int("maxAttempts", maxAttempts))

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Where("industry_naics_code IS NULL OR industry_naics_code = ''").
		Where("industry_request_count IS NULL OR industry_request_count < ?", maxAttempts).
		Where("industry_requested_at IS NULL OR industry_requested_at < ?", utils.Now().Add(-1*time.Hour*time.Duration(hoursFromPreviousAttempt))).
		Order("CASE WHEN industry_requested_at IS NULL THEN 0 ELSE 1 END ASC").
		Order("CASE WHEN industry_requested_at IS NULL THEN created_at END DESC").
		Order("industry_requested_at ASC").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogFields(tracingLog.Int("found", len(organizations)))
	return organizations, nil
}

func (r *globalOrganizationRepository) GetOrganizationsToEnrichDescription(ctx context.Context, hoursFromPreviousAttempt, maxAttempts, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.GetOrganizationsToEnrichDescription")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Int("hoursFromPreviousAttempt", hoursFromPreviousAttempt), tracingLog.Int("limit", limit), tracingLog.Int("maxAttempts", maxAttempts))

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Where("description IS NULL OR description = ''").
		Where("description_request_count IS NULL OR description_request_count < ?", maxAttempts).
		Where("description_requested_at IS NULL OR description_requested_at < ?", utils.Now().Add(-1*time.Hour*time.Duration(hoursFromPreviousAttempt))).
		Order("CASE WHEN description_requested_at IS NULL THEN 0 ELSE 1 END ASC").
		Order("CASE WHEN description_requested_at IS NULL THEN created_at END DESC").
		Order("description_requested_at ASC").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogFields(tracingLog.Int("found", len(organizations)))
	return organizations, nil
}

func (r *globalOrganizationRepository) GetOrganizationsToEnrichName(ctx context.Context, hoursFromPreviousAttempt, maxAttempts, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.GetOrganizationsToEnrichName")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Int("hoursFromPreviousAttempt", hoursFromPreviousAttempt),
		tracingLog.Int("limit", limit),
		tracingLog.Int("maxAttempts", maxAttempts))

	// Condition to flag “suspicious” or missing name
	// 1) Name is NULL or empty
	// 2) Name is all-caps (letters/numbers/spaces)
	// 3) Name has excessive punctuation (example threshold: >2 punctuation chars)
	// 4) Name looks like a domain
	// 5) Name is a known placeholder (“none”, “n/a”, “na”, “unknown”, “null”, etc.)
	suspiciousNameCondition := `
        (
            name IS NULL
            OR name = ''
            OR name ~ '^[A-Z0-9\\s]+$'
            OR length(regexp_replace(name, '[a-zA-Z0-9\\s]', '', 'g')) > 2
            OR name ~ '^[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$'
            OR lower(name) IN ('none','n/a','na','unknown','null')
        )
    `

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Where(suspiciousNameCondition).
		Where("name_request_count IS NULL OR name_request_count < ?", maxAttempts).
		Where("name_requested_at IS NULL OR name_requested_at < ?", utils.Now().Add(-1*time.Hour*time.Duration(hoursFromPreviousAttempt))).
		Order("CASE WHEN name_requested_at IS NULL THEN 0 ELSE 1 END ASC").
		Order("CASE WHEN name_requested_at IS NULL THEN created_at END DESC").
		Order("name_requested_at ASC").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogFields(tracingLog.Int("found", len(organizations)))
	return organizations, nil
}

func (r *globalOrganizationRepository) MarkIndustryEnrichRequested(ctx context.Context, id uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.MarkIndustryEnrichRequested")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Uint64("id", id))

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		UpdateColumn("industry_request_count", gorm.Expr("COALESCE(industry_request_count, 0) + 1")).
		UpdateColumn("industry_requested_at", utils.Now())
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return result.Error
	}
	return nil
}

func (r *globalOrganizationRepository) MarkDescriptionEnrichRequested(ctx context.Context, id uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.MarkDescriptionEnrichRequested")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Uint64("id", id))

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		UpdateColumn("description_request_count", gorm.Expr("COALESCE(description_request_count, 0) + 1")).
		UpdateColumn("description_requested_at", utils.Now())
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return result.Error
	}
	return nil
}

func (r *globalOrganizationRepository) MarkNameEnrichRequested(ctx context.Context, id uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.MarkNameEnrichRequested")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Uint64("id", id))

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		UpdateColumn("name_request_count", gorm.Expr("COALESCE(name_request_count, 0) + 1")).
		UpdateColumn("name_requested_at", utils.Now())
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return result.Error
	}
	return nil
}

func (r *globalOrganizationRepository) SetIndustry(ctx context.Context, id uint64, industryNaicsCode, industryNaicsName string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.SetIndustry")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Uint64("id", id), tracingLog.String("industryNaicsCode", industryNaicsCode), tracingLog.String("industryNaicsName", industryNaicsName))

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"industry_naics_code": industryNaicsCode,
			"industry_naics_name": industryNaicsName,
			"industry_set_at":     utils.Now(),
		})
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return result.Error
	}
	return nil
}

func (r *globalOrganizationRepository) SetDescription(ctx context.Context, id uint64, description string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.SetDescription")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Uint64("id", id), tracingLog.String("description", description))

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"description":        description,
			"description_set_at": utils.Now(),
		})
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return result.Error
	}
	return nil
}

func (r *globalOrganizationRepository) SetName(ctx context.Context, id uint64, name string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.SetName")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Uint64("id", id), tracingLog.String("name", name))

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"name_set_at": utils.Now(),
		})
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return result.Error
	}
	return nil
}

func (r *globalOrganizationRepository) GetGlobalOrganizationsToSyncIntoTenantOrganizations(ctx context.Context, daysFromPreviousSync, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.GetGlobalOrganizationsToSyncIntoTenantOrganizations")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Int("daysFromPreviousSync", daysFromPreviousSync), tracingLog.Int("limit", limit))

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Where("((industry_naics_code IS NOT NULL AND industry_naics_code <> '') AND (description IS NOT NULL AND description <> ''))").
		Where("synced_to_neo_at IS NULL OR synced_to_neo_at < ?", utils.Now().Add(-24*time.Hour*time.Duration(daysFromPreviousSync))).
		Order("CASE WHEN synced_to_neo_at IS NULL THEN 0 ELSE 1 END ASC").
		Order("CASE WHEN synced_to_neo_at IS NULL THEN created_at END DESC").
		Order("synced_to_neo_at ASC").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogFields(tracingLog.Int("found", len(organizations)))
	return organizations, nil
}

func (r *globalOrganizationRepository) MarkGlobalOrganizationSyncedToNeo(ctx context.Context, id uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationRepository.MarkGlobalOrganizationSyncedToNeo")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Uint64("id", id))

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		UpdateColumn("synced_to_neo_at", utils.Now())
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return result.Error
	}
	return nil
}
