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

	query := `
        INSERT INTO cache_email_validation_domain (
            domain, is_catch_all, is_firewalled, can_connect_smtp, provider, firewall, created_at, updated_at
        ) VALUES (
            ?, ?, ?, ?, ?, ?, ?, ?
        )
        ON CONFLICT (domain) DO UPDATE SET
            is_catch_all = EXCLUDED.is_catch_all,
            is_firewalled = EXCLUDED.is_firewalled,
            can_connect_smtp = EXCLUDED.can_connect_smtp,
            provider = EXCLUDED.provider,
            firewall = EXCLUDED.firewall,
            updated_at = EXCLUDED.updated_at
        RETURNING *
    `

	now := utils.Now()
	var result entity.CacheEmailValidationDomain

	err := r.db.WithContext(ctx).Raw(query,
		cacheEmailValidationDomain.Domain,
		cacheEmailValidationDomain.IsCatchAll,
		cacheEmailValidationDomain.IsFirewalled,
		cacheEmailValidationDomain.CanConnectSMTP,
		cacheEmailValidationDomain.Provider,
		cacheEmailValidationDomain.Firewall,
		now,
		now,
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}
