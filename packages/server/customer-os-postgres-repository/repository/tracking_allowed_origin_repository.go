package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type TrackingAllowedOriginRepository interface {
	GetTenantForOrigin(ctx context.Context, origin string) (string, error)
	Create(ctx context.Context, whitelist postgres_entity.TrackingAllowedOrigin) (*postgres_entity.TrackingAllowedOrigin, error)
	Update(ctx context.Context, whitelist postgres_entity.TrackingAllowedOrigin) (*postgres_entity.TrackingAllowedOrigin, error)
	FindAll(ctx context.Context, whitelist postgres_entity.TrackingAllowedOrigin) (*[]postgres_entity.TrackingAllowedOrigin, error)
	Find(ctx context.Context, whitelist postgres_entity.TrackingAllowedOrigin) (*postgres_entity.TrackingAllowedOrigin, error)
}

type trackingAllowedOriginRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewTrackingAllowedOriginRepository(gormDb *gorm.DB) TrackingAllowedOriginRepository {
	return &trackingAllowedOriginRepositoryImpl{gormDb: gormDb}
}

func (repo *trackingAllowedOriginRepositoryImpl) Create(ctx context.Context, whitelist postgres_entity.TrackingAllowedOrigin) (*postgres_entity.TrackingAllowedOrigin, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingAllowedOriginRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var created postgres_entity.TrackingAllowedOrigin
	err := repo.gormDb.Create(&whitelist).Scan(&created).Error
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (repo *trackingAllowedOriginRepositoryImpl) FindAll(ctx context.Context, whitelist postgres_entity.TrackingAllowedOrigin) (*[]postgres_entity.TrackingAllowedOrigin, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingAllowedOriginRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var results []postgres_entity.TrackingAllowedOrigin
	err := repo.gormDb.
		Where(&whitelist).
		Find(&results).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &results, nil
}

func (repo *trackingAllowedOriginRepositoryImpl) Find(ctx context.Context, whitelist postgres_entity.TrackingAllowedOrigin) (*postgres_entity.TrackingAllowedOrigin, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingAllowedOriginRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var result postgres_entity.TrackingAllowedOrigin
	if err := repo.gormDb.Where(&whitelist).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &result, nil
}

func (repo *trackingAllowedOriginRepositoryImpl) Update(ctx context.Context, whitelist postgres_entity.TrackingAllowedOrigin) (*postgres_entity.TrackingAllowedOrigin, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TrackingAllowedOriginRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var updated postgres_entity.TrackingAllowedOrigin
	err := repo.gormDb.Model(&whitelist).
		Updates(&whitelist).
		First(&updated).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updated, nil
}

func (repo *trackingAllowedOriginRepositoryImpl) GetTenantForOrigin(ctx context.Context, origin string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingAllowedOriginRepository.GetTenantForOrigin")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("origin", origin)

	var result postgres_entity.TrackingAllowedOrigin
	err := repo.gormDb.Model(&postgres_entity.TrackingAllowedOrigin{}).
		Where("origin = ? OR origin = ?", origin, origin+"/").
		Order("created_at ASC").
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		tracing.TraceErr(span, err)
		return "", err
	}

	return result.Tenant, nil
}
