package postgres_repository

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type TenantRepository interface {
	Create(ctx context.Context, tenantEntity postgres_entity.Tenant) (*postgres_entity.Tenant, error)
	GetTenantByHashId(ctx context.Context, hashID string) (string, error)
	GetHashID(ctx context.Context, tenant string) (string, error)
	PermanentlyDelete(ctx context.Context, tenant string) error
	SetCompanyReport(ctx context.Context, companyReport string) error
}

type tenantRepository struct {
	gormDb *gorm.DB
}

func NewTenantRepository(gormDb *gorm.DB) TenantRepository {
	return &tenantRepository{gormDb: gormDb}
}

func (e *tenantRepository) GetTenantByHashId(ctx context.Context, hashID string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantRepository.GetTenantByHashId")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("hashID", hashID)

	var tenant postgres_entity.Tenant
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

	var tenant postgres_entity.Tenant
	err := e.gormDb.
		Where("name = ?", tenantName).
		First(&tenant).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "failed to get tenant")
	}

	return tenant.HashID, nil
}

func (e *tenantRepository) Create(ctx context.Context, tenantEntity postgres_entity.Tenant) (*postgres_entity.Tenant, error) {
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

	err := e.gormDb.Where("name = ?", tenant).Delete(&postgres_entity.Tenant{}).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "failed to delete tenant")
	}
	return nil
}

func (e *tenantRepository) SetCompanyReport(ctx context.Context, companyReport string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantRepository.SetCompanyReport")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		return coserrors.ErrTenantNotSet
	}
	result := e.gormDb.Model(&postgres_entity.Tenant{}).
		Where("name = ?", tenant).
		Update("company_report", companyReport)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("Tenant not found")
	}

	return nil
}
