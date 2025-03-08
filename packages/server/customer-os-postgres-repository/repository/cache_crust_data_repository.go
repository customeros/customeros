package postgres_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type cacheCrustDataRepository struct {
	db *gorm.DB
}

type CacheCrustDataRepository interface {
	// Create creates a new cache entry
	Create(ctx context.Context, data postgres_entity.CacheCrustData) (*postgres_entity.CacheCrustData, error)

	// GetByCompanyAndTitle retrieves a cache entry by company domain and job titles
	GetByCompanyAndTitle(ctx context.Context, company string, titles []string) (*postgres_entity.CacheCrustData, error)
}

func NewCacheCrustDataRepository(db *gorm.DB) CacheCrustDataRepository {
	return &cacheCrustDataRepository{
		db: db,
	}
}

func (r *cacheCrustDataRepository) Create(ctx context.Context, data postgres_entity.CacheCrustData) (*postgres_entity.CacheCrustData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CrustDataCacheRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create crust data cache"))
		return nil, err
	}

	return &data, nil
}

func (r *cacheCrustDataRepository) GetByCompanyAndTitle(ctx context.Context, company string, titles []string) (*postgres_entity.CacheCrustData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CrustDataCacheRepository.GetByCompanyAndTitle")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var result postgres_entity.CacheCrustData
	if err := r.db.WithContext(ctx).
		Where("request_company = ? AND request_title = ?", company, pq.Array(titles)).
		First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		tracing.TraceErr(span, errors.Wrap(err, "failed to get crust data cache"))
		return nil, err
	}

	return &result, nil
}
