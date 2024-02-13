package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	"gorm.io/gorm"
)

type ApiKeyRepo struct {
	db *gorm.DB
}

type ApiKeyRepository interface {
	GetDefaultApiKey(tenant string) helper.QueryResult
	GetApiKey(tenant, name string) helper.QueryResult
	GetApiKeys(tenant string) helper.QueryResult
	SaveApiKey(integration entity.ApiKey) helper.QueryResult
}

func NewApiKeyRepo(db *gorm.DB) *ApiKeyRepo {
	return &ApiKeyRepo{db: db}
}

func (r *ApiKeyRepo) GetDefaultApiKey(tenant string) helper.QueryResult {
	var apiKeyEntity entity.ApiKey
	err := r.db.
		Where("tenant_name = ?", tenant).
		Where("name = ?", "default").
		First(&apiKeyEntity).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &apiKeyEntity}
}

func (r *ApiKeyRepo) GetApiKey(tenant, name string) helper.QueryResult {
	var apiKeyEntity entity.ApiKey
	err := r.db.
		Where("tenant_name = ?", tenant).
		Where("name = ?", name).
		First(&apiKeyEntity).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &apiKeyEntity}
}

func (r *ApiKeyRepo) GetApiKeys(tenant string) helper.QueryResult {
	var apiKeyEntities []entity.ApiKey
	err := r.db.
		Where("tenant_name = ?", tenant).
		Find(&apiKeyEntities).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &apiKeyEntities}
}

func (r *ApiKeyRepo) SaveApiKey(apiKey entity.ApiKey) helper.QueryResult {
	apiKeyEntity := entity.ApiKey{
		TenantName: apiKey.TenantName,
		Name:       apiKey.Name,
		Key:        apiKey.Key,
	}

	err := r.db.Create(&apiKeyEntity).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &apiKeyEntity}
}
