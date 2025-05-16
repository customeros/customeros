package postgres_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
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
	spans, _ := telemetry.StartPostgresSpan(ctx, "TenantRepository.GetTenantByHashId")
	defer spans.Finish()

	spans.LogKV("hashID", hashID)

	var tenant postgres_entity.Tenant
	err := e.gormDb.
		Where("tenant_hash = ?", hashID).
		First(&tenant).Error
	if err != nil {
		spans.TraceError(err)
		return "", errors.Wrap(err, "failed to get tenant")
	}

	return tenant.Name, nil
}

func (e *tenantRepository) GetHashID(ctx context.Context, tenantName string) (string, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TenantRepository.GetHashID")
	defer spans.Finish()

	var tenant postgres_entity.Tenant
	err := e.gormDb.
		Where("name = ?", tenantName).
		First(&tenant).Error
	if err != nil {
		spans.TraceError(err)
		return "", errors.Wrap(err, "failed to get tenant")
	}

	return tenant.HashID, nil
}

func (e *tenantRepository) Create(ctx context.Context, tenantEntity postgres_entity.Tenant) (*postgres_entity.Tenant, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TenantRepository.Create")
	defer spans.Finish()

	if tenantEntity.Name == "" {
		return nil, fmt.Errorf("no tenant name, cannot create tenant")
	}

	tenantEntity.HashID = utils.GenerateHashId(tenantEntity.Name, 12)

	err := e.gormDb.Create(&tenantEntity).Error
	if err != nil {
		spans.TraceError(err)
		return nil, errors.Wrap(err, "failed to create tenant")
	}

	return &tenantEntity, nil
}

func (e *tenantRepository) PermanentlyDelete(ctx context.Context, tenant string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "TenantRepository.PermanentlyDelete")
	defer spans.Finish()

	err := e.gormDb.Where("name = ?", tenant).Delete(&postgres_entity.Tenant{}).Error
	if err != nil {
		spans.TraceError(err)
		return errors.Wrap(err, "failed to delete tenant")
	}
	return nil
}

func (e *tenantRepository) SetCompanyReport(ctx context.Context, companyReport string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TenantRepository.SetCompanyReport")
	defer spans.Finish()

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
