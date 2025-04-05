package postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type GlobalOrganizationRepository interface {
	GetById(ctx context.Context, id uint64) (*postgres_entity.GlobalOrganization, error)
	GetByPrimaryDomain(ctx context.Context, domain string) (*postgres_entity.GlobalOrganization, error)
	GetByPrimaryDomains(ctx context.Context, domains []string) ([]*postgres_entity.GlobalOrganization, error)
	Create(ctx context.Context, organization *postgres_entity.GlobalOrganization) (*postgres_entity.GlobalOrganization, error)
	CreateIfNotExists(ctx context.Context, primaryDomain string) (*postgres_entity.GlobalOrganization, error)
	Update(ctx context.Context, organization *postgres_entity.GlobalOrganization) (*postgres_entity.GlobalOrganization, error)
	Delete(ctx context.Context, id uint64) error
	Search(ctx context.Context, searchTerm string, limit int) ([]*postgres_entity.GlobalOrganization, error)
	GetOrganizationsToEnrichIndustry(ctx context.Context, hoursFromPreviousAttempt, maxAttempts, limit int) ([]*postgres_entity.GlobalOrganization, error)
	GetOrganizationsToEnrichDescription(ctx context.Context, hoursFromPreviousAttempt, maxAttempts, limit int) ([]*postgres_entity.GlobalOrganization, error)
	GetOrganizationsToEnrichName(ctx context.Context, hoursFromPreviousAttempt, maxAttempts, limit int) ([]*postgres_entity.GlobalOrganization, error)
	GetOrganizationsToScrape(ctx context.Context, limit int) ([]*postgres_entity.GlobalOrganization, error)
	GetOrganizationsToFetchLogo(ctx context.Context, limit int) ([]*postgres_entity.GlobalOrganization, error)
	GetOrganizationsToFetchIcon(ctx context.Context, limit int) ([]*postgres_entity.GlobalOrganization, error)
	MarkIndustryEnrichRequested(ctx context.Context, id uint64) error
	MarkDescriptionEnrichRequested(ctx context.Context, id uint64) error
	MarkNameEnrichRequested(ctx context.Context, id uint64) error
	SetScrapeStatus(ctx context.Context, id uint64, status enum.ScrapeStatus) error
	SetIndustry(ctx context.Context, id uint64, industryNaicsCode, industryNaicsName string) error
	SetDescription(ctx context.Context, id uint64, description string) error
	SetName(ctx context.Context, id uint64, name string) error
	SetLogo(ctx context.Context, id uint64, logoPath string) error
	SetIcon(ctx context.Context, id uint64, iconPath string) error
	SetDownloadStatusLogo(ctx context.Context, id uint64, status enum.DownloadStatus) error
	SetDownloadStatusIcon(ctx context.Context, id uint64, status enum.DownloadStatus) error
	GetGlobalOrganizationsToSyncIntoTenantOrganizations(ctx context.Context, daysFromPreviousSync, limit int) ([]*postgres_entity.GlobalOrganization, error)
	MarkGlobalOrganizationSyncedToNeo(ctx context.Context, id uint64) error
	GetByLinkedInUrl(ctx context.Context, linkedInUrl string) (*postgres_entity.GlobalOrganization, error)
	AddOtherSocials(ctx context.Context, primaryDomain string, otherSocials []string) error
}

type globalOrganizationRepository struct {
	db *gorm.DB
}

func NewGlobalOrganizationRepository(gormDb *gorm.DB) GlobalOrganizationRepository {
	return &globalOrganizationRepository{db: gormDb}
}

func (r *globalOrganizationRepository) GetById(ctx context.Context, id uint64) (*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.GetById")
	defer spans.Finish()
	spans.LogKV("id", id)

	organization := &postgres_entity.GlobalOrganization{}
	result := r.db.WithContext(ctx).Where("id = ?", id).First(organization)
	if result.Error != nil {
		spans.LogKV("result.found", false)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("result.found", true)
	return organization, nil
}

func (r *globalOrganizationRepository) GetByPrimaryDomains(ctx context.Context, domains []string) ([]*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.GetByPrimaryDomains")
	defer spans.Finish()
	spans.LogKV("domains", domains)

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).Where("primary_domain IN ?", domains).Find(&organizations)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("result.count", len(organizations))
	return organizations, nil
}

func (r *globalOrganizationRepository) GetByPrimaryDomain(ctx context.Context, domain string) (*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.GetByPrimaryDomain")
	defer spans.Finish()
	spans.LogKV("domain", domain)

	organization := &postgres_entity.GlobalOrganization{}
	result := r.db.WithContext(ctx).Where("primary_domain = ?", domain).First(organization)
	if result.Error != nil {
		spans.LogKV("result.found", false)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("result.found", true)
	return organization, nil
}

func (r *globalOrganizationRepository) Create(ctx context.Context, organization *postgres_entity.GlobalOrganization) (*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.Create")
	defer spans.Finish()
	spans.LogObjectAsJson("organization", organization)

	result := r.db.WithContext(ctx).Create(&organization)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	return organization, nil
}

func (r *globalOrganizationRepository) CreateIfNotExists(ctx context.Context, primaryDomain string) (*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.CreateIfNotExists")
	defer spans.Finish()
	spans.LogKV("primaryDomain", primaryDomain)

	// Check if organization exists
	existing, err := r.GetByPrimaryDomain(ctx, primaryDomain)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if existing != nil {
		spans.LogKV("result.created", false)
		return existing, nil
	}

	// Create new organization
	organization := &postgres_entity.GlobalOrganization{
		PrimaryDomain: primaryDomain,
		CreatedAt:     utils.Now(),
		UpdatedAt:     utils.Now(),
	}

	result := r.db.WithContext(ctx).Create(organization)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}

	spans.LogKV("result.created", true)
	return organization, nil
}

func (r *globalOrganizationRepository) Update(ctx context.Context, organization *postgres_entity.GlobalOrganization) (*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.Update")
	defer spans.Finish()
	spans.LogObjectAsJson("organization", organization)

	organization.UpdatedAt = utils.Now()
	result := r.db.WithContext(ctx).Save(organization)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	return organization, nil
}

func (r *globalOrganizationRepository) Search(ctx context.Context, searchTerm string, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.SearchOrganizations")
	defer spans.Finish()
	spans.LogKV("searchTerm", searchTerm)

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Select("id, name, primary_domain, website, logo_url, icon_url, other_domains").
		Where("name ILIKE ? OR primary_domain ILIKE ? OR other_domains ILIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("found", len(organizations))
	return organizations, nil
}

func (r *globalOrganizationRepository) GetOrganizationsToEnrichIndustry(ctx context.Context, hoursFromPreviousAttempt, maxAttempts, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.GetOrganizationsToEnrichIndustry")
	defer spans.Finish()
	spans.LogKV("hoursFromPreviousAttempt", hoursFromPreviousAttempt, "limit", limit, "maxAttempts", maxAttempts)

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Where("active = ?", true).
		Where("industry_naics_code IS NULL OR industry_naics_code = ''").
		Where("industry_request_count IS NULL OR industry_request_count < ?", maxAttempts).
		Where("industry_requested_at IS NULL OR industry_requested_at < ?", utils.Now().Add(-1*time.Hour*time.Duration(hoursFromPreviousAttempt))).
		Order("CASE WHEN industry_requested_at IS NULL THEN 0 ELSE 1 END ASC").
		Order("CASE WHEN industry_requested_at IS NULL THEN created_at END DESC").
		Order("industry_requested_at ASC").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("found", len(organizations))
	return organizations, nil
}

func (r *globalOrganizationRepository) GetOrganizationsToScrape(ctx context.Context, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.GetOrganizationsToScrape")
	defer spans.Finish()
	spans.LogKV("limit", limit)

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Where("(scrape_status = ? OR scrape_status IS NULL) OR (scrape_status = ? AND (scraped_at < ? OR scraped_at IS NULL))",
			enum.ScrapeNotScraped.String(),
			enum.ScrapeError.String(),
			utils.Now().Add(-48*time.Hour)).
		Order("created_at DESC").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("found", len(organizations))
	return organizations, nil
}

func (r *globalOrganizationRepository) GetOrganizationsToFetchLogo(ctx context.Context, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.GetOrganizationsToFetchLogo")
	defer spans.Finish()
	spans.LogKV("limit", limit)

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Where("download_status_logo = ? OR download_status_logo IS NULL", enum.DownloadNotStarted.String()).
		Where("logo_url IS NOT NULL AND logo_url != ''").
		Order("created_at DESC").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("result.count", len(organizations))
	return organizations, nil
}

func (r *globalOrganizationRepository) GetOrganizationsToFetchIcon(ctx context.Context, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.GetOrganizationsToFetchIcon")
	defer spans.Finish()
	spans.LogKV("limit", limit)

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Where("download_status_icon = ? OR download_status_icon IS NULL", enum.DownloadNotStarted.String()).
		Where("icon_url IS NOT NULL AND icon_url != ''").
		Order("created_at DESC").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("result.count", len(organizations))
	return organizations, nil
}

func (r *globalOrganizationRepository) GetOrganizationsToEnrichDescription(ctx context.Context, hoursFromPreviousAttempt, maxAttempts, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.GetOrganizationsToEnrichDescription")
	defer spans.Finish()
	spans.LogKV("hoursFromPreviousAttempt", hoursFromPreviousAttempt, "limit", limit, "maxAttempts", maxAttempts)

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Where("active = ?", true).
		Where("description IS NULL OR description = ''").
		Where("description_request_count IS NULL OR description_request_count < ?", maxAttempts).
		Where("description_requested_at IS NULL OR description_requested_at < ?", utils.Now().Add(-1*time.Hour*time.Duration(hoursFromPreviousAttempt))).
		Order("CASE WHEN description_requested_at IS NULL THEN 0 ELSE 1 END ASC").
		Order("CASE WHEN description_requested_at IS NULL THEN created_at END DESC").
		Order("description_requested_at ASC").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("found", len(organizations))
	return organizations, nil
}

func (r *globalOrganizationRepository) GetOrganizationsToEnrichName(ctx context.Context, hoursFromPreviousAttempt, maxAttempts, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.GetOrganizationsToEnrichName")
	defer spans.Finish()
	spans.LogKV("hoursFromPreviousAttempt", hoursFromPreviousAttempt, "limit", limit, "maxAttempts", maxAttempts)

	// Condition to flag "suspicious" or missing name
	// 1) Name is NULL or empty
	// 2) Name is all-caps (letters/numbers/spaces)
	// 3) Name has excessive punctuation (example threshold: >2 punctuation chars)
	// 4) Name looks like a domain
	// 5) Name is a known placeholder ("none", "n/a", "na", "unknown", "null", etc.)
	// 6) Name is all lowercase (letters/numbers/spaces)
	// 7) Name has a known suffix (" ltd", " inc", " gmbh", " s.a", " plc", " pty", " corp", " co", " sa")
	// 8) Name has a punctuation character
	suspiciousNameCondition := `
        (
            name IS NULL
            OR name = ''
            OR name ~ '^[A-Z0-9\s]+$'
            OR length(regexp_replace(name, '[a-zA-Z0-9\s]', '', 'g')) > 2
            OR name ~ '^[a-zA-Z0-9.-]+\.[a-zA-Z0-9.-]+$'
            OR lower(name) IN ('none','n/a','na','unknown','null')
			OR name ~ '^[a-z0-9\s]+$'
			OR lower(name) ~ '( ltd| inc| gmbh| plc| pty| corp| co| llc| llp| ag| bv| as| ab| nv| se| sl| sc| cv| sa| sarl| spa| srl| srls| snc| sas)'
			OR name ~ '.*[.,;:!?].*'
        )
    `

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Where(suspiciousNameCondition).
		Where("active = ?", true).
		Where("name_request_count IS NULL OR name_request_count < ?", maxAttempts).
		Where("name_requested_at IS NULL OR name_requested_at < ?", utils.Now().Add(-1*time.Hour*time.Duration(hoursFromPreviousAttempt))).
		Order("CASE WHEN name_requested_at IS NULL THEN 0 ELSE 1 END ASC").
		Order("CASE WHEN name_requested_at IS NULL THEN created_at END DESC").
		Order("name_requested_at ASC").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("found", len(organizations))
	return organizations, nil
}

func (r *globalOrganizationRepository) MarkIndustryEnrichRequested(ctx context.Context, id uint64) error {
	return r.markEnrichRequested(ctx, id, "industry")
}

func (r *globalOrganizationRepository) MarkDescriptionEnrichRequested(ctx context.Context, id uint64) error {
	return r.markEnrichRequested(ctx, id, "description")
}

func (r *globalOrganizationRepository) MarkNameEnrichRequested(ctx context.Context, id uint64) error {
	return r.markEnrichRequested(ctx, id, "name")
}

func (r *globalOrganizationRepository) SetIndustry(ctx context.Context, id uint64, industryNaicsCode, industryNaicsName string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.SetIndustry")
	defer spans.Finish()
	spans.LogKV("id", id, "industryNaicsCode", industryNaicsCode, "industryNaicsName", industryNaicsName)

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"industry_naics_code": industryNaicsCode,
			"industry_naics_name": industryNaicsName,
			"industry_set_at":     utils.Now(),
		})
	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}
	return nil
}

func (r *globalOrganizationRepository) SetDescription(ctx context.Context, id uint64, description string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.SetDescription")
	defer spans.Finish()
	spans.LogKV("id", id, "description", description)

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"description":        description,
			"description_set_at": utils.Now(),
		})
	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}
	return nil
}

func (r *globalOrganizationRepository) SetName(ctx context.Context, id uint64, name string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.SetName")
	defer spans.Finish()
	spans.LogKV("id", id, "name", name)

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"name_set_at": utils.Now(),
		})
	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}
	return nil
}

func (r *globalOrganizationRepository) GetGlobalOrganizationsToSyncIntoTenantOrganizations(ctx context.Context, daysFromPreviousSync, limit int) ([]*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.GetGlobalOrganizationsToSyncIntoTenantOrganizations")
	defer spans.Finish()
	spans.LogKV("daysFromPreviousSync", daysFromPreviousSync, "limit", limit)

	organizations := make([]*postgres_entity.GlobalOrganization, 0)
	result := r.db.WithContext(ctx).
		Where("active = ?", true).
		Where("(synced_to_neo_at IS NULL OR (synced_to_neo_at < ? AND updated_at > synced_to_neo_at))", utils.Now().Add(-24*time.Hour*time.Duration(daysFromPreviousSync))).
		Order("synced_to_neo_at IS NULL DESC, COALESCE(synced_to_neo_at, created_at) ASC").
		Limit(limit).
		Find(&organizations)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("found", len(organizations))
	return organizations, nil
}

func (r *globalOrganizationRepository) MarkGlobalOrganizationSyncedToNeo(ctx context.Context, id uint64) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.MarkGlobalOrganizationSyncedToNeo")
	defer spans.Finish()
	spans.LogKV("id", id)

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		UpdateColumn("synced_to_neo_at", utils.Now().Add(10*time.Second))
	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}
	return nil
}

func (r *globalOrganizationRepository) markEnrichRequested(ctx context.Context, id uint64, fieldBase string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.markEnrichRequested")
	defer spans.Finish()
	spans.LogKV("id", id, "fieldBase", fieldBase)

	countField := fieldBase + "_request_count" // e.g., "industry_request_count"
	timeField := fieldBase + "_requested_at"   // e.g., "industry_requested_at"

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		UpdateColumn(countField, gorm.Expr("COALESCE("+countField+", 0) + 1")).
		UpdateColumn(timeField, utils.Now())
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *globalOrganizationRepository) SetScrapeStatus(ctx context.Context, id uint64, status enum.ScrapeStatus) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.SetScrapeStatus")
	defer spans.Finish()
	spans.LogKV("id", id, "status", status.String())

	active := status == enum.ScrapeCompleted
	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"scrape_status":  status.String(),
			"scrape_attempt": gorm.Expr("COALESCE(scrape_attempt, 0) + 1"),
			"scraped_at":     utils.Now(),
			"active":         active,
		})

	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		err := errors.New("no organization found with provided ID")
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *globalOrganizationRepository) SetLogo(ctx context.Context, id uint64, logoPath string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.SetLogo")
	defer spans.Finish()
	spans.LogKV("id", id, "logoPath", logoPath)

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"logo_path":       logoPath,
			"download_status": enum.DownloadCompleted.String(),
		})

	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		err := errors.New("no organization found with provided ID")
		spans.TraceError(err)
		return err
	}

	spans.LogKV("rows_affected", result.RowsAffected)
	return nil
}

func (r *globalOrganizationRepository) SetIcon(ctx context.Context, id uint64, iconPath string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.SetIcon")
	defer spans.Finish()
	spans.LogKV("id", id, "iconPath", iconPath)

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"icon_path":       iconPath,
			"download_status": enum.DownloadCompleted.String(),
		})

	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		err := errors.New("no organization found with provided ID")
		spans.TraceError(err)
		return err
	}

	spans.LogKV("rows_affected", result.RowsAffected)
	return nil
}

func (r *globalOrganizationRepository) SetDownloadStatusLogo(ctx context.Context, id uint64, status enum.DownloadStatus) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.SetDownloadStatusLogo")
	defer spans.Finish()
	spans.LogKV("id", id, "status", status.String())

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"download_status_logo": status.String(),
		})

	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}

	return nil
}

func (r *globalOrganizationRepository) SetDownloadStatusIcon(ctx context.Context, id uint64, status enum.DownloadStatus) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.SetDownloadStatusIcon")
	defer spans.Finish()
	spans.LogKV("id", id, "status", status.String())

	result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"download_status_icon": status.String(),
		})

	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}

	return nil
}

func (r *globalOrganizationRepository) GetByLinkedInUrl(ctx context.Context, linkedInUrl string) (*postgres_entity.GlobalOrganization, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.GetByLinkedInUrl")
	defer spans.Finish()
	spans.LogKV("linkedInUrl", linkedInUrl)

	organization := &postgres_entity.GlobalOrganization{}
	result := r.db.WithContext(ctx).Where("linkedin = ?", linkedInUrl).First(organization)
	if result.Error != nil {
		spans.LogKV("found", false)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	spans.LogKV("found", true)
	return organization, nil
}

func (r *globalOrganizationRepository) AddOtherSocials(ctx context.Context, primaryDomain string, otherSocials []string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.AddOtherSocials")
	defer spans.Finish()
	spans.LogKV("primaryDomain", primaryDomain, "otherSocials", otherSocials)

	// Get existing organization or create new one
	org, err := r.CreateIfNotExists(ctx, primaryDomain)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	// Add new socials if they don't exist
	added := false
	for _, social := range otherSocials {
		social = strings.TrimSpace(social)
		if social == "" {
			continue
		}
		exists := false
		for _, existing := range org.OtherSocials {
			if existing == social {
				exists = true
				break
			}
		}
		if !exists {
			org.OtherSocials = append(org.OtherSocials, social)
			added = true
		}
	}

	// Update if new socials were added
	if added {
		result := r.db.WithContext(ctx).Model(&postgres_entity.GlobalOrganization{}).
			Where("id = ?", org.ID).
			Updates(map[string]interface{}{
				"other_socials": org.OtherSocials,
			})
		if result.Error != nil {
			spans.TraceError(result.Error)
			return result.Error
		}
	}

	return nil
}

func (r *globalOrganizationRepository) Delete(ctx context.Context, id uint64) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationRepository.Delete")
	defer spans.Finish()
	spans.TagEntity(fmt.Sprintf("%d", id))

	result := r.db.WithContext(ctx).Delete(&postgres_entity.GlobalOrganization{}, id)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return result.Error
	}

	return nil
}
