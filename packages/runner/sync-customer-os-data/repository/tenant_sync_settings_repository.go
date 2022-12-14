package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"gorm.io/gorm"
)

type TenantSyncSettingsRepository interface {
	GetTenantsForSync() (entity.TenantSyncSettingsList, error)
}

type tenantSyncSettingsRepository struct {
	db *gorm.DB
}

func NewTenantSyncSettingsRepository(gormDb *gorm.DB) TenantSyncSettingsRepository {
	return &tenantSyncSettingsRepository{
		db: gormDb,
	}
}

func (r *tenantSyncSettingsRepository) GetTenantsForSync() (entity.TenantSyncSettingsList, error) {
	var tenantsForSync entity.TenantSyncSettingsList

	err := r.db.
		Where(&entity.TenantSyncSettings{Enabled: true}).
		Find(&tenantsForSync).Error
	if err != nil {
		return nil, err
	}

	return tenantsForSync, nil
}
