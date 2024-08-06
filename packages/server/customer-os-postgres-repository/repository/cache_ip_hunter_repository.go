package repository

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type CacheIpHunterRepository interface {
	Get(ctx context.Context, ip string) (*entity.CacheIpHunter, error)
	Save(ctx context.Context, cacheIpHunter entity.CacheIpHunter) (*entity.CacheIpHunter, error)
}

type cacheIpHunterRepository struct {
	db *gorm.DB
}

func NewCacheIpHunterRepository(gormDb *gorm.DB) CacheIpHunterRepository {
	return &cacheIpHunterRepository{db: gormDb}
}

func (r cacheIpHunterRepository) Get(ctx context.Context, ip string) (*entity.CacheIpHunter, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheIpHunterRepository.Register")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("ip", ip))

	var cacheIpHunter entity.CacheIpHunter
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

func (r cacheIpHunterRepository) Save(ctx context.Context, cacheIpHunter entity.CacheIpHunter) (*entity.CacheIpHunter, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheIpHunterRepository.Register")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "cacheIpHunter", cacheIpHunter)

	var existingHunter entity.CacheIpHunter
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
			return nil, err
		}
		cacheIpHunter = existingHunter
	}

	return &cacheIpHunter, nil
}
