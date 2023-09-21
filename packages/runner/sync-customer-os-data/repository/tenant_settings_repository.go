package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TenantSettingsRepository interface {
	GetTenantSettings(ctx context.Context, tenant string) (entity.TenantSettings, error)
}

type tenantSettingsRepository struct {
	db *gorm.DB
}

func NewTenantSettingsRepository(gormDb *gorm.DB) TenantSettingsRepository {
	return &tenantSettingsRepository{
		db: gormDb,
	}
}

func (r *tenantSettingsRepository) GetTenantSettings(ctx context.Context, tenant string) (entity.TenantSettings, error) {
	var tenantSettings entity.TenantSettings
	result := r.db.Where("tenant_name = ?", tenant).First(&tenantSettings)
	if result.Error != nil {
		return tenantSettings, result.Error
	}
	return tenantSettings, nil
}
