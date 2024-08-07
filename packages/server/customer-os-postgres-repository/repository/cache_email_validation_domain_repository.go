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

type CacheEmailValidationDomainRepository interface {
	Get(ctx context.Context, domain string) (*entity.CacheEmailValidationDomain, error)
	Save(ctx context.Context, cacheEmailValidationDomain entity.CacheEmailValidationDomain) (*entity.CacheEmailValidationDomain, error)
}

type cacheEmailValidationDomainRepository struct {
	db *gorm.DB
}

func NewCacheEmailValidationDomainRepository(gormDb *gorm.DB) CacheEmailValidationDomainRepository {
	return &cacheEmailValidationDomainRepository{db: gormDb}
}

func (r cacheEmailValidationDomainRepository) Get(ctx context.Context, domain string) (*entity.CacheEmailValidationDomain, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailValidationDomainRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("domain", domain))

	var cacheEmailValidationDomain entity.CacheEmailValidationDomain
	result := r.db.WithContext(ctx).Where("domain = ?", domain).First(&cacheEmailValidationDomain)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	return &cacheEmailValidationDomain, nil
}

func (r cacheEmailValidationDomainRepository) Save(ctx context.Context, cacheEmailValidationDomain entity.CacheEmailValidationDomain) (*entity.CacheEmailValidationDomain, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailValidationDomainRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "cacheEmailValidationDomain", cacheEmailValidationDomain)

	var existingData entity.CacheEmailValidationDomain
	result := r.db.WithContext(ctx).Where("domain = ?", cacheEmailValidationDomain.Domain).First(&existingData)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Record doesn't exist, create a new one
			cacheEmailValidationDomain.CreatedAt = utils.Now()
			cacheEmailValidationDomain.UpdatedAt = utils.Now()
			if err := r.db.WithContext(ctx).Create(&cacheEmailValidationDomain).Error; err != nil {
				return nil, err
			}
		} else {
			// Some other error occurred
			return nil, result.Error
		}
	} else {
		// Record exists, update it
		updates := map[string]interface{}{
			"updated_at":       utils.Now(),
			"is_catch_all":     cacheEmailValidationDomain.IsCatchAll,
			"is_firewalled":    cacheEmailValidationDomain.IsFirewalled,
			"can_connect_smtp": cacheEmailValidationDomain.CanConnectSMTP,
			"provider":         cacheEmailValidationDomain.Provider,
			"firewall":         cacheEmailValidationDomain.Firewall,
		}
		if err := r.db.WithContext(ctx).Model(&existingData).Updates(updates).Error; err != nil {
			return nil, err
		}
		cacheEmailValidationDomain = existingData
	}

	return &cacheEmailValidationDomain, nil
}
