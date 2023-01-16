package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/helper"
	"gorm.io/gorm"
)

type TenantSettingsRepository interface {
	FindForTenantName(tenantId string) helper.QueryResult
	Save(tenantSettings entity.TenantSettings) helper.QueryResult
}

type tenantSettingsRepo struct {
	db *gorm.DB
}

func NewTenantSettingsRepository(db *gorm.DB) TenantSettingsRepository {
	return &tenantSettingsRepo{
		db: db,
	}
}

func (r *tenantSettingsRepo) FindForTenantName(tenantId string) helper.QueryResult {
	var tenantSettings entity.TenantSettings

	err := r.db.
		Where("tenant_id = ?", tenantId).
		First(&tenantSettings).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return helper.QueryResult{Error: err}
	}
	if err == gorm.ErrRecordNotFound {
		return helper.QueryResult{Result: nil}
	}

	return helper.QueryResult{Result: tenantSettings}
}

func (r *tenantSettingsRepo) Save(tenantSettings entity.TenantSettings) helper.QueryResult {

	result := r.db.Save(&tenantSettings)

	if result.Error != nil {
		return helper.QueryResult{Error: result.Error}
	}

	return helper.QueryResult{Result: tenantSettings}
}
