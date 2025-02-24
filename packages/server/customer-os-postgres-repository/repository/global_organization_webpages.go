package postgres_repository

import (
	"context"
	"errors"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type GlobalOrganizationWebpageRepository interface {
	Save(ctx context.Context, webpageData postgres_entity.GlobalOrganizationWebpages) (*postgres_entity.GlobalOrganizationWebpages, error)
	GetWebpage(ctx context.Context, url string, lookbackInDays int) (*postgres_entity.GlobalOrganizationWebpages, error)
}

type globalOrganizationWebpageRepository struct {
	gormDb *gorm.DB
}

func NewGlobalOrganizationWebpageRepository(gormDb *gorm.DB) GlobalOrganizationWebpageRepository {
	return &globalOrganizationWebpageRepository{gormDb: gormDb}
}

func (r *globalOrganizationWebpageRepository) Save(ctx context.Context, webpageData postgres_entity.GlobalOrganizationWebpages) (*postgres_entity.GlobalOrganizationWebpages, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "globalOrganizationWebpageRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// Check if a record with this URL already exists
	var existingRecord postgres_entity.GlobalOrganizationWebpages
	result := r.gormDb.Where("url = ?", webpageData.Url).First(&existingRecord)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// URL doesn't exist, create a new record
			webpageData.CreatedAt = time.Now()
			webpageData.UpdatedAt = time.Now()
			if err := r.gormDb.Create(&webpageData).Error; err != nil {
				return nil, err
			}
			return &webpageData, nil
		}
		// Some other database error occurred
		return nil, result.Error
	}

	// URL exists, update the content and updated_at timestamp
	existingRecord.Content = webpageData.Content
	existingRecord.UpdatedAt = time.Now()

	// If primary domain is provided, update it too
	if webpageData.PrimaryDomain != "" {
		existingRecord.PrimaryDomain = webpageData.PrimaryDomain
	}

	if err := r.gormDb.Save(&existingRecord).Error; err != nil {
		return nil, err
	}

	return &existingRecord, nil
}

func (r *globalOrganizationWebpageRepository) GetWebpage(ctx context.Context, url string, lookbackInDays int) (*postgres_entity.GlobalOrganizationWebpages, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "globalOrganizationWebpageRepository.GetWebpage")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// Calculate the cutoff date
	cutoffDate := time.Now().AddDate(0, 0, -lookbackInDays)

	var record postgres_entity.GlobalOrganizationWebpages
	err := r.gormDb.
		Where("url = ?", url).
		Where("updated_at > ?", cutoffDate).
		First(&record).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &record, nil
}
