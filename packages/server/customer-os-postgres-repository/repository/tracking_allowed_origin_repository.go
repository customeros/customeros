package repository

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

type TrackingAllowedOriginRepository interface {
	GetTenantForOrigin(ctx context.Context, origin string) (*string, error)
}

type trackingAllowedOriginRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewTrackingAllowedOriginRepository(gormDb *gorm.DB) TrackingAllowedOriginRepository {
	return &trackingAllowedOriginRepositoryImpl{gormDb: gormDb}
}

func (repo *trackingAllowedOriginRepositoryImpl) GetTenantForOrigin(ctx context.Context, origin string) (*string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingAllowedOriginRepository.GetTenantForOrigin")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("origin", origin)

	var result entity.TrackingAllowedOrigin
	err := repo.gormDb.Model(&entity.TrackingAllowedOrigin{}).
		Where("origin = ? OR origin = ?", origin, origin+"/").
		Order("created_at ASC").
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	if result.Tenant == "" {
		return nil, nil
	}

	return &result.Tenant, nil
}
