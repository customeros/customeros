package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository/helper"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type PostmarkApiKeyRepo struct {
	db *gorm.DB
}

type PostmarkApiKeyRepository interface {
	GetPostmarkApiKey(ctx context.Context, tenant string) helper.QueryResult
	CreateApiKey(ctx context.Context, integration postgres_entity.PostmarkApiKey) helper.QueryResult
}

func NewPostmarkApiKeyRepo(db *gorm.DB) *PostmarkApiKeyRepo {
	return &PostmarkApiKeyRepo{db: db}
}

func (r *PostmarkApiKeyRepo) GetPostmarkApiKey(ctx context.Context, tenant string) helper.QueryResult {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "PostmarkApiKeyRepo.GetPostmarkApiKey")
	defer spans.Finish()

	var postmarkApiKeyEntity postgres_entity.PostmarkApiKey
	err := r.db.
		Where("tenant_name = ?", tenant).
		First(&postmarkApiKeyEntity).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &postmarkApiKeyEntity}
}

func (r *PostmarkApiKeyRepo) CreateApiKey(ctx context.Context, apiKey postgres_entity.PostmarkApiKey) helper.QueryResult {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "PostmarkApiKeyRepo.CreateApiKey")
	defer spans.Finish()

	postmarkApiKeyEntity := postgres_entity.PostmarkApiKey{
		TenantName: apiKey.TenantName,
		Key:        apiKey.Key,
	}

	err := r.db.Create(&postmarkApiKeyEntity).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &postmarkApiKeyEntity}
}
