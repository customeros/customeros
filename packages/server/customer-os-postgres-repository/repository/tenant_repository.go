package repository

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type TenantRepository interface {
	Create(ctx context.Context, tenantEntity entity.Tenant) (*entity.Tenant, error)
	GetTenant(ctx context.Context, hashID string) (string, error)
	GetHashID(ctx context.Context, tenant string) (string, error)
	PermanentlyDelete(ctx context.Context, tenant string) error
}

type tenantRepository struct {
	gormDb *gorm.DB
}

func NewTenantRepository(gormDb *gorm.DB) TenantRepository {
	return &tenantRepository{gormDb: gormDb}
}

func (e *tenantRepository) GetTenant(ctx context.Context, hashID string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantRepository.GetTenant")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var tenant entity.Tenant
	err := e.gormDb.
		Where("tenant_hash = ?", hashID).
		First(&tenant).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "failed to get tenant")
	}

	return tenant.Name, nil
}

func (e *tenantRepository) GetHashID(ctx context.Context, tenantName string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantRepository.GetHashID")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var tenant entity.Tenant
	err := e.gormDb.
		Where("name = ?", tenantName).
		First(&tenant).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "failed to get tenant")
	}

	return tenant.HashID, nil
}

func (e *tenantRepository) Create(ctx context.Context, tenantEntity entity.Tenant) (*entity.Tenant, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if tenantEntity.Name == "" {
		return nil, fmt.Errorf("No tenant name, cannot create tenant")
	}

	tenantEntity.HashID = utils.GenerateHashId(tenantEntity.Name, 12)

	err := e.gormDb.Create(&tenantEntity).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "failed to create tenant")
	}

	return &tenantEntity, nil
}

func (e *tenantRepository) PermanentlyDelete(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.PermanentlyDelete")
	defer span.Finish()

	err := e.gormDb.Where("name = ?", tenant).Delete(&entity.Tenant{}).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "failed to delete tenant")
	}
	return nil
}
