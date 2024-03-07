package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	"gorm.io/gorm"
)

type TenantWebhookApiKeyRepo struct {
	db *gorm.DB
}

type TenantWebhookApiKeyRepository interface {
	GetTenantWithApiKey(apiKey string) helper.QueryResult
	SaveApiKey(integration entity.TenantWebhookApiKey) helper.QueryResult
}

func NewTenantWebhookApiKeyRepo(db *gorm.DB) *TenantWebhookApiKeyRepo {
	return &TenantWebhookApiKeyRepo{db: db}
}

func (r *TenantWebhookApiKeyRepo) GetTenantWithApiKey(apiKey string) helper.QueryResult {
	var apiKeyEntity entity.TenantWebhookApiKey
	err := r.db.
		Where("key = ?", apiKey).
		First(&apiKeyEntity).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &apiKeyEntity}
}

func (r *TenantWebhookApiKeyRepo) SaveApiKey(apiKey entity.TenantWebhookApiKey) helper.QueryResult {
	apiKeyEntity := entity.TenantWebhookApiKey{
		TenantName: apiKey.TenantName,
		Key:        apiKey.Key,
	}

	err := r.db.Create(&apiKeyEntity).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &apiKeyEntity}
}
