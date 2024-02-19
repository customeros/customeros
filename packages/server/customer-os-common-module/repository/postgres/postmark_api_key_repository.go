package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	"gorm.io/gorm"
)

type PostmarkApiKeyRepo struct {
	db *gorm.DB
}

type PostmarkApiKeyRepository interface {
	GetPostmarkApiKey(tenant string) helper.QueryResult
	CreateApiKey(integration entity.PostmarkApiKey) helper.QueryResult
}

func NewPostmarkApiKeyRepo(db *gorm.DB) *PostmarkApiKeyRepo {
	return &PostmarkApiKeyRepo{db: db}
}

func (r *PostmarkApiKeyRepo) GetPostmarkApiKey(tenant string) helper.QueryResult {
	var postmarkApiKeyEntity entity.PostmarkApiKey
	err := r.db.
		Where("tenant_name = ?", tenant).
		First(&postmarkApiKeyEntity).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &postmarkApiKeyEntity}
}

func (r *PostmarkApiKeyRepo) CreateApiKey(apiKey entity.PostmarkApiKey) helper.QueryResult {
	postmarkApiKeyEntity := entity.PostmarkApiKey{
		TenantName: apiKey.TenantName,
		Key:        apiKey.Key,
	}

	err := r.db.Create(&postmarkApiKeyEntity).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &postmarkApiKeyEntity}
}
