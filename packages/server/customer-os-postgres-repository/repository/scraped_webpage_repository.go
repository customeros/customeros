package postgres_repository

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go/log"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
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
	SetContentStage(ctx context.Context, url string, contentStage enum.CustomerJourneyStage) error
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "scrapedWebpageRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ScrapedWebpageRepository.GetWebpage")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogKV("url", url, "lookbackInDays", lookbackInDays)

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
			span.LogFields(log.Bool("result.found", false))
			return nil, nil
		}
		return nil, err
	}

	tracing.LogObjectAsJson(span, "result.webpage", record)
	return &record, nil
}

func (r *scrapedWebpageRepository) GetAllWebpagesByPrimaryDomains(ctx context.Context, primaryDomains []string) ([]*postgres_entity.ScrapedWebpage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "scrapedWebpageRepository.GetAllWebpagesByPrimaryDomains")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "scrapedWebpageRepository.SetLinks")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.Model(&postgres_entity.ScrapedWebpage{}).
		Where("url = ?", url).
		Update("links", pq.StringArray(links)).
		Error

	return err
}

func (r *scrapedWebpageRepository) SetContent(ctx context.Context, url string, content string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "scrapedWebpageRepository.SetContent")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.Model(&postgres_entity.ScrapedWebpage{}).
		Where("url = ?", url).
		Update("content", content).
		Error

	return err
}

func (r *scrapedWebpageRepository) SetContentStage(ctx context.Context, url string, contentStage enum.CustomerJourneyStage) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "scrapedWebpageRepository.SetContentStage")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.Model(&postgres_entity.ScrapedWebpage{}).
		Where("url = ?", url).
		Update("content_stage", contentStage).
		Error

	return err
}

func (r *scrapedWebpageRepository) SetWebpageTopics(ctx context.Context, url string, topics []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "scrapedWebpageRepository.SetWebpageTopics")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.Model(&postgres_entity.ScrapedWebpage{}).
		Where("url = ?", url).
		Update("topics", pq.StringArray(topics)).
		Error

	return err
}

func (r *scrapedWebpageRepository) SetWebpageCategory(ctx context.Context, url string, category enum.WebpageCategory) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "scrapedWebpageRepository.SetWebpageCategory")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.Model(&postgres_entity.ScrapedWebpage{}).
		Where("url = ?", url).
		Update("category", category.String()).
		Error

	return err
}

func (r *scrapedWebpageRepository) GetWebpagesWithoutLinks(ctx context.Context, limit int) ([]*postgres_entity.ScrapedWebpage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "scrapedWebpageRepository.GetWebpagesWithoutLinks")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	webpages := make([]*postgres_entity.ScrapedWebpage, 0)
	result := r.gormDb.WithContext(ctx).
		Where("array_length(links, 1) IS NULL OR array_length(links, 1) = 0").
		Order("created_at DESC").
		Limit(limit).
		Find(&webpages)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	span.LogKV("found", len(webpages))
	return webpages, nil
}
