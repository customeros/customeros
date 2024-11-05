package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TenantRepository interface {
	Create(ctx context.Context, tenantEntity entity.Tenant) (*entity.Tenant, error)
	PermanentlyDelete(ctx context.Context, tenant string) error
}

type tenantRepository struct {
	gormDb *gorm.DB
}

func (e tenantRepository) PermanentlyDelete(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.PermanentlyDelete")
	defer span.Finish()

	err := e.gormDb.Where("name = ?", tenant).Delete(&entity.Tenant{}).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "failed to delete tenant")
	}
	return nil
}

func NewTenantRepository(gormDb *gorm.DB) TenantRepository {
	return &tenantRepository{gormDb: gormDb}
}

func (e tenantRepository) Create(ctx context.Context, tenantEntity entity.Tenant) (*entity.Tenant, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := e.gormDb.Create(&tenantEntity).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "failed to create tenant")
	}

	return &tenantEntity, nil
}
