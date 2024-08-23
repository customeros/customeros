package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository/helper"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type PostmarkApiKeyRepo struct {
	db *gorm.DB
}

type PostmarkApiKeyRepository interface {
	GetPostmarkApiKey(ctx context.Context, tenant string) helper.QueryResult
	CreateApiKey(ctx context.Context, integration entity.PostmarkApiKey) helper.QueryResult
}

func NewPostmarkApiKeyRepo(db *gorm.DB) *PostmarkApiKeyRepo {
	return &PostmarkApiKeyRepo{db: db}
}

func (r *PostmarkApiKeyRepo) GetPostmarkApiKey(ctx context.Context, tenant string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "PostmarkApiKeyRepo.GetPostmarkApiKey")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	var postmarkApiKeyEntity entity.PostmarkApiKey
	err := r.db.
		Where("tenant_name = ?", tenant).
		First(&postmarkApiKeyEntity).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &postmarkApiKeyEntity}
}

func (r *PostmarkApiKeyRepo) CreateApiKey(ctx context.Context, apiKey entity.PostmarkApiKey) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "PostmarkApiKeyRepo.CreateApiKey")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

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
