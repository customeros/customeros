package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	"gorm.io/gorm"
)

type TenantApiKeyRepo struct {
	db *gorm.DB
}

type TenantApiKeyRepository interface {
	GetApiKey(tenant string) helper.QueryResult
	SaveApiKey(integration entity.TenantApiKey) helper.QueryResult
}

func NewTenantApiKeyRepo(db *gorm.DB) *TenantApiKeyRepo {
	return &TenantApiKeyRepo{db: db}
}

func (r *TenantApiKeyRepo) GetApiKey(tenant string) helper.QueryResult {
	var apiKeyEntity entity.TenantApiKey
	err := r.db.
		Where("tenant_name = ?", tenant).
		First(&apiKeyEntity).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &apiKeyEntity}
}

func (r *TenantApiKeyRepo) SaveApiKey(apiKey entity.TenantApiKey) helper.QueryResult {
	apiKeyEntity := entity.TenantApiKey{
		TenantName: apiKey.TenantName,
		Key:        apiKey.Key,
	}

	err := r.db.Create(&apiKeyEntity).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &apiKeyEntity}
}
