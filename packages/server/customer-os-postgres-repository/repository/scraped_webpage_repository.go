package postgres_repository

import (
	"context"
	"errors"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type ScrapedWebpageRepository interface {
	Save(ctx context.Context, webpageData postgres_entity.ScrapedWebpage) (*postgres_entity.ScrapedWebpage, error)
	GetWebpage(ctx context.Context, url string, lookbackInDays int) (*postgres_entity.ScrapedWebpage, error)
	GetAllWebpagesByPrimaryDomains(ctx context.Context, primaryDomains []string) ([]*postgres_entity.ScrapedWebpage, error)
	GetWebpagesWithoutLinks(ctx context.Context, limit int) ([]*postgres_entity.ScrapedWebpage, error)
	SetLinks(ctx context.Context, url string, links []string) error
	SetContent(ctx context.Context, url string, content string) error
	SetWebpageTopics(ctx context.Context, url string, topics []string) error
	SetWebpageCategory(ctx context.Context, url string, category enum.WebpageCategory) error
}

type scrapedWebpageRepository struct {
	gormDb *gorm.DB
}

func NewScrapedWebpageRepository(gormDb *gorm.DB) ScrapedWebpageRepository {
	return &scrapedWebpageRepository{gormDb: gormDb}
}

func (r *scrapedWebpageRepository) Save(ctx context.Context, webpageData postgres_entity.ScrapedWebpage) (*postgres_entity.ScrapedWebpage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "scrapedWebpageRepository.Save")
	defer spans.Finish()

	// Check if a record with this URL already exists
	var existingRecord postgres_entity.ScrapedWebpage
	result := r.gormDb.Where("url = ?", webpageData.Url).First(&existingRecord)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// URL doesn't exist, create a new record
			webpageData.CreatedAt = time.Now()
			webpageData.UpdatedAt = time.Now()

			// Make sure Links is properly handled for PostgreSQL array
			if err := r.gormDb.Create(&webpageData).Error; err != nil {
				return nil, err
			}
			return &webpageData, nil
		}
		// Some other database error occurred
		return nil, result.Error
	}

	// Update existing record fields if they're provided
	if webpageData.Content != "" {
		existingRecord.Content = webpageData.Content
		existingRecord.UpdatedAt = time.Now()
	}

	if webpageData.PrimaryDomain != "" {
		existingRecord.PrimaryDomain = webpageData.PrimaryDomain
	}

	// Add handling for Links field
	if len(webpageData.Links) > 0 {
		existingRecord.Links = webpageData.Links
		existingRecord.UpdatedAt = time.Now()
	}

	if err := r.gormDb.Save(&existingRecord).Error; err != nil {
		return nil, err
	}

	return &existingRecord, nil
}

func (r *scrapedWebpageRepository) GetWebpage(ctx context.Context, url string, lookbackInDays int) (*postgres_entity.ScrapedWebpage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "ScrapedWebpageRepository.GetWebpage")
	defer spans.Finish()
	spans.LogKV("url", url, "lookbackInDays", lookbackInDays)

	// Calculate the cutoff date
	cutoffDate := utils.Now().AddDate(0, 0, -lookbackInDays)

	var record postgres_entity.ScrapedWebpage
	err := r.gormDb.
		Where("url = ?", url).
		Where("updated_at > ?", cutoffDate).
		First(&record).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		return nil, err
	}

	spans.LogObjectAsJson("result.webpage", record)
	return &record, nil
}

func (r *scrapedWebpageRepository) GetAllWebpagesByPrimaryDomains(ctx context.Context, primaryDomains []string) ([]*postgres_entity.ScrapedWebpage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "scrapedWebpageRepository.GetAllWebpagesByPrimaryDomains")
	defer spans.Finish()
	spans.LogObjectAsJson("primaryDomains", primaryDomains)

	var records []*postgres_entity.ScrapedWebpage

	err := r.gormDb.
		Where("primary_domain IN ?", primaryDomains).
		Find(&records).
		Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *scrapedWebpageRepository) SetLinks(ctx context.Context, url string, links []string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "ScrapedWebpageRepository.SetLinks")
	defer spans.Finish()
	spans.LogKV("url", url)
	spans.LogObjectAsJson("links", links)

	err := r.gormDb.Model(&postgres_entity.ScrapedWebpage{}).
		Where("url = ?", url).
		Updates(map[string]interface{}{
			"links":            pq.StringArray(links),
			"links_checked_at": time.Now(),
		}).
		Error

	return err
}

func (r *scrapedWebpageRepository) SetContent(ctx context.Context, url string, content string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "ScrapedWebpageRepository.SetContent")
	defer spans.Finish()
	spans.LogKV("url", url, "content.length", len(content))

	err := r.gormDb.Model(&postgres_entity.ScrapedWebpage{}).
		Where("url = ?", url).
		Update("content", content).
		Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *scrapedWebpageRepository) SetWebpageTopics(ctx context.Context, url string, topics []string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "ScrapedWebpageRepository.SetWebpageTopics")
	defer spans.Finish()
	spans.LogKV("url", url, "topics.length", len(topics))
	spans.LogObjectAsJson("topics", topics)

	err := r.gormDb.Model(&postgres_entity.ScrapedWebpage{}).
		Where("url = ?", url).
		Update("topics", pq.StringArray(topics)).
		Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *scrapedWebpageRepository) SetWebpageCategory(ctx context.Context, url string, category enum.WebpageCategory) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "ScrapedWebpageRepository.SetWebpageCategory")
	defer spans.Finish()
	spans.LogKV("url", url, "category", category)

	err := r.gormDb.Model(&postgres_entity.ScrapedWebpage{}).
		Where("url = ?", url).
		Update("category", category.String()).
		Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

// Number of days to recheck links
const LinksCheckStaleThresholdDays = 30

func (r *scrapedWebpageRepository) GetWebpagesWithoutLinks(ctx context.Context, limit int) ([]*postgres_entity.ScrapedWebpage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "ScrapedWebpageRepository.GetWebpagesWithoutLinks")
	defer spans.Finish()

	webpages := make([]*postgres_entity.ScrapedWebpage, 0)

	// Calculate half of the limit
	halfLimit := limit / 2
	if halfLimit == 0 {
		halfLimit = 1
	}

	// Calculate the cutoff date for stale links
	staleDate := utils.Now().AddDate(0, 0, -LinksCheckStaleThresholdDays)

	// Get recent records
	var recentWebpages []*postgres_entity.ScrapedWebpage
	resultRecent := r.gormDb.WithContext(ctx).
		Where("array_length(links, 1) IS NULL OR array_length(links, 1) = 0").
		Where("links_checked_at IS NULL OR links_checked_at < ?", staleDate).
		Order("created_at DESC").
		Limit(halfLimit).
		Find(&recentWebpages)
	if resultRecent.Error != nil {
		spans.TraceError(resultRecent.Error)
		return nil, resultRecent.Error
	}

	// Track URLs we've already seen to avoid duplicates
	seenUrls := make(map[string]bool)
	for _, webpage := range recentWebpages {
		seenUrls[webpage.Url] = true
		webpages = append(webpages, webpage)
	}

	// Get random records, excluding URLs we already have
	remainingLimit := limit - len(webpages)
	if remainingLimit > 0 {
		var randomWebpages []*postgres_entity.ScrapedWebpage
		excludeUrls := make([]string, 0, len(seenUrls))
		for url := range seenUrls {
			excludeUrls = append(excludeUrls, url)
		}

		resultRandom := r.gormDb.WithContext(ctx).
			Where("array_length(links, 1) IS NULL OR array_length(links, 1) = 0").
			Where("links_checked_at IS NULL OR links_checked_at < ?", staleDate).
			Where("url NOT IN (?)", excludeUrls).
			Order("RANDOM()").
			Limit(remainingLimit).
			Find(&randomWebpages)
		if resultRandom.Error != nil {
			spans.TraceError(resultRandom.Error)
			return nil, resultRandom.Error
		}

		webpages = append(webpages, randomWebpages...)
	}

	spans.LogKV("result.count", len(webpages))
	return webpages, nil
}
