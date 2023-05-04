package repository

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/helper"
	"gorm.io/gorm"
)

type TenantSettingsRepository interface {
	FindForTenantName(tenantName string) helper.QueryResult
	Save(tenantSettings *entity.TenantSettings) helper.QueryResult
	SaveKey(keys []entity.TenantAPIKey) error
	CheckKeysExist(tenantName string, keyName string) (bool, error)
}

type tenantSettingsRepo struct {
	db *gorm.DB
}

func NewTenantSettingsRepository(db *gorm.DB) TenantSettingsRepository {
	return &tenantSettingsRepo{
		db: db,
	}
}

func (r *tenantSettingsRepo) FindForTenantName(tenantName string) helper.QueryResult {
	var tenantSettings entity.TenantSettings

	err := r.db.
		Where("tenant_name = ?", tenantName).
		First(&tenantSettings).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return helper.QueryResult{Error: err}
	}
	if err == gorm.ErrRecordNotFound {
		return helper.QueryResult{Result: nil}
	}

	return helper.QueryResult{Result: tenantSettings}
}

func (r *tenantSettingsRepo) CheckKeysExist(tenantName string, keyName string) (bool, error) {
	var rows int64
	err := r.db.
		Where(&entity.TenantAPIKey{TenantName: tenantName, Key: keyName}, "tenant_name", "key").Count(&rows).Error

	if err != nil {
		return false, fmt.Errorf("CheckKeysExist: %w", err)
	}

	if rows == 0 {
		return false, nil
	}
	return true, nil
}

func (r *tenantSettingsRepo) SaveKey(keys []entity.TenantAPIKey) error {

	for _, key := range keys {
		result := r.db.Save(key)
		if result.Error != nil {
			return fmt.Errorf("SaveKey: %w", result.Error)
		}
	}
	return nil
}

func (r *tenantSettingsRepo) Save(tenantSettings *entity.TenantSettings) helper.QueryResult {

	result := r.db.Save(tenantSettings)

	if result.Error != nil {
		return helper.QueryResult{Error: result.Error}
	}

	return helper.QueryResult{Result: tenantSettings}
}
