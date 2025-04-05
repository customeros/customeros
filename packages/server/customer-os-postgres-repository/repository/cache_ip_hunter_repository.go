package postgres_repository

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type CacheIpHunterRepository interface {
	Get(ctx context.Context, ip string) (*postgres_entity.CacheIpHunter, error)
	Save(ctx context.Context, cacheIpHunter postgres_entity.CacheIpHunter) (*postgres_entity.CacheIpHunter, error)
}

type cacheIpHunterRepository struct {
	db *gorm.DB
}

func NewCacheIpHunterRepository(gormDb *gorm.DB) CacheIpHunterRepository {
	return &cacheIpHunterRepository{db: gormDb}
}

func (r cacheIpHunterRepository) Get(ctx context.Context, ip string) (*postgres_entity.CacheIpHunter, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheIpHunterRepository.Register")
	defer spans.Finish()
	spans.LogKV("ip", ip)

	var cacheIpHunter postgres_entity.CacheIpHunter
	result := r.db.WithContext(ctx).Where("ip = ?", ip).First(&cacheIpHunter)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	return &cacheIpHunter, nil
}

func (r cacheIpHunterRepository) Save(ctx context.Context, cacheIpHunter postgres_entity.CacheIpHunter) (*postgres_entity.CacheIpHunter, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CacheIpHunterRepository.Register")
	defer spans.Finish()
	spans.LogObjectAsJson("cacheIpHunter", cacheIpHunter)

	var existingHunter postgres_entity.CacheIpHunter
	result := r.db.WithContext(ctx).Where("ip = ?", cacheIpHunter.Ip).First(&existingHunter)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Record doesn't exist, create a new one
			cacheIpHunter.CreatedAt = utils.Now()
			cacheIpHunter.UpdatedAt = utils.Now()
			if err := r.db.WithContext(ctx).Create(&cacheIpHunter).Error; err != nil {
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
			"data":       cacheIpHunter.Data,
		}
		if err := r.db.WithContext(ctx).Model(&existingHunter).Updates(updates).Error; err != nil {
			spans.TraceError(err)
			return nil, err
		}
		cacheIpHunter = existingHunter
	}

	return &cacheIpHunter, nil
}
