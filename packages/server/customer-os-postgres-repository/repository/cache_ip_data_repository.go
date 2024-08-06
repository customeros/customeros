package repository

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
	"time"
)

type CacheIpDataRepository interface {
	Get(ctx context.Context, ip string) (*entity.CacheIpData, error)
	Save(ctx context.Context, cacheIpData entity.CacheIpData) (*entity.CacheIpData, error)
}

type cacheIpDataRepository struct {
	db *gorm.DB
}

func NewCacheIpDataRepository(gormDb *gorm.DB) CacheIpDataRepository {
	return &cacheIpDataRepository{db: gormDb}
}

func (r cacheIpDataRepository) Get(ctx context.Context, ip string) (*entity.CacheIpData, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheIpDataRepository.Register")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("ip", ip))

	var cacheIpData entity.CacheIpData
	result := r.db.WithContext(ctx).Where("ip = ?", ip).First(&cacheIpData)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	return &cacheIpData, nil
}

func (r cacheIpDataRepository) Save(ctx context.Context, cacheIpData entity.CacheIpData) (*entity.CacheIpData, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheIpDataRepository.Register")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "cacheIpData", cacheIpData)

	var existingData entity.CacheIpData
	result := r.db.WithContext(ctx).Where("ip = ?", cacheIpData.Ip).First(&existingData)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Record doesn't exist, create a new one
			cacheIpData.CreatedAt = time.Now()
			cacheIpData.UpdatedAt = time.Now()
			if err := r.db.WithContext(ctx).Create(&cacheIpData).Error; err != nil {
				return nil, err
			}
		} else {
			// Some other error occurred
			return nil, result.Error
		}
	} else {
		// Record exists, update it
		updates := map[string]interface{}{
			"updated_at": time.Now(),
		}
		if err := r.db.WithContext(ctx).Model(&existingData).Updates(updates).Error; err != nil {
			return nil, err
		}
		cacheIpData = existingData
	}

	return &cacheIpData, nil
}
