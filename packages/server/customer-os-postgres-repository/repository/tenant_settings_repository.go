package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TenantSettingsRepository interface {
	FindForTenantName(ctx context.Context, tenantName string) (*postgres_entity.TenantSettings, error)
	Save(ctx context.Context, tenantSettings *postgres_entity.TenantSettings) (*postgres_entity.TenantSettings, error)
}

type tenantSettingsRepo struct {
	db *gorm.DB
}

func NewTenantSettingsRepository(db *gorm.DB) TenantSettingsRepository {
	return &tenantSettingsRepo{
		db: db,
	}
}

func (r *tenantSettingsRepo) FindForTenantName(ctx context.Context, tenantName string) (*postgres_entity.TenantSettings, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TenantSettingsRepository.FindForTenantName")
	defer spans.Finish()

	var tenantSettings postgres_entity.TenantSettings

	err := r.db.
		Where("tenant_name = ?", tenantName).
		First(&tenantSettings).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		spans.TraceError(err)
		return nil, err
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		spans.LogKV("result.found", false)
		return nil, nil
	}

	spans.LogKV("result.found", true)

	return &tenantSettings, nil
}

func (r *tenantSettingsRepo) Save(ctx context.Context, tenantSettings *postgres_entity.TenantSettings) (*postgres_entity.TenantSettings, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TenantSettingsRepository.Save")
	defer spans.Finish()

	err := r.db.Save(tenantSettings).Error

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return tenantSettings, nil
}
