package postgres_repository

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type CacheIpDataRepository interface {
	Get(ctx context.Context, ip string) (*postgres_entity.CacheIpData, error)
	Save(ctx context.Context, cacheIpData postgres_entity.CacheIpData) (*postgres_entity.CacheIpData, error)
}

type cacheIpDataRepository struct {
	db *gorm.DB
}

func NewCacheIpDataRepository(gormDb *gorm.DB) CacheIpDataRepository {
	return &cacheIpDataRepository{db: gormDb}
}

func (r cacheIpDataRepository) Get(ctx context.Context, ip string) (*postgres_entity.CacheIpData, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheIpDataRepository.Get")
	defer spans.Finish()
	spans.LogKV("ip", ip)

	var cacheIpData postgres_entity.CacheIpData
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

func (r cacheIpDataRepository) Save(ctx context.Context, cacheIpData postgres_entity.CacheIpData) (*postgres_entity.CacheIpData, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheIpDataRepository.Save")
	defer spans.Finish()
	spans.LogObjectAsJson("cacheIpData", cacheIpData)

	var existingData postgres_entity.CacheIpData
	result := r.db.WithContext(ctx).Where("ip = ?", cacheIpData.Ip).First(&existingData)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Record doesn't exist, create a new one
			cacheIpData.CreatedAt = utils.Now()
			cacheIpData.UpdatedAt = utils.Now()
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
			"updated_at": utils.Now(),
			"data":       cacheIpData.Data,
		}
		if err := r.db.WithContext(ctx).Model(&existingData).Updates(updates).Error; err != nil {
			spans.TraceError(err)
			return nil, err
		}
		cacheIpData = existingData
	}

	return &cacheIpData, nil
}
