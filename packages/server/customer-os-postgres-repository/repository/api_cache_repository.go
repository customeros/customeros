package repository

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type ApiCacheRepository interface {
	GetAll(ctx context.Context) ([]*entity.ApiCache, error)
	Get(ctx context.Context, tenant, typee string) (*entity.ApiCache, error)
	Save(ctx context.Context, entity entity.ApiCache) error
}

type apiCacheRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewApiCacheRepository(gormDb *gorm.DB) ApiCacheRepository {
	return &apiCacheRepositoryImpl{gormDb: gormDb}
}

func (r *apiCacheRepositoryImpl) GetAll(ctx context.Context) ([]*entity.ApiCache, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ApiCacheRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var entities []*entity.ApiCache
	err := r.gormDb.Find(&entities).Error

	return entities, err
}

func (r *apiCacheRepositoryImpl) Get(ctx context.Context, tenant, typee string) (*entity.ApiCache, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ApiCacheRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("typee", typee))

	var entity entity.ApiCache
	err := r.gormDb.Where("tenant = ? AND type = ?", tenant, typee).First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &entity, err
}

func (r *apiCacheRepositoryImpl) Save(ctx context.Context, entity entity.ApiCache) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ApiCacheRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, entity.Tenant)

	if entity.ID != "" {
		return r.gormDb.Create(&entity).Error
	} else {
		existing, err := r.Get(ctx, entity.Tenant, entity.Type)
		if err != nil {
			return err
		}

		if existing != nil {
			entity.ID = existing.ID
			return r.gormDb.Save(&entity).Error
		} else {
			return r.gormDb.Create(&entity).Error
		}
	}
}
