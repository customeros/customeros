package repository

import (
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TenantWebhookApiKeyRepository interface {
	CreateApiKey(ctx context.Context, tenant string) error
	GetTenantForApiKey(ctx context.Context, apiKey string) (*entity.TenantWebhookApiKey, error)
	GetFirstApiKeyForTenant(ctx context.Context, tenant string) (*entity.TenantWebhookApiKey, error)
}

type tenantWebhookApiKeyRepository struct {
	gormDb *gorm.DB
}

func NewTenantWebhookApiKeyRepository(gormDb *gorm.DB) TenantWebhookApiKeyRepository {
	return &tenantWebhookApiKeyRepository{gormDb: gormDb}
}

func (r *tenantWebhookApiKeyRepository) CreateApiKey(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWebhookApiKeyRepository.CreateApiKey")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	now := utils.Now()
	apiKey := entity.TenantWebhookApiKey{
		Tenant:    tenant,
		Key:       entity.KeyPrefix + utils.GenerateKey(32),
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := r.gormDb.Create(&apiKey).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *tenantWebhookApiKeyRepository) GetTenantForApiKey(ctx context.Context, apiKey string) (*entity.TenantWebhookApiKey, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWebhookApiKeyRepository.GetTenantWithApiKey")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// get record for api key or nil if not found
	var apiKeyRecord entity.TenantWebhookApiKey
	err := r.gormDb.
		Where("key = ?", apiKey).
		First(&apiKeyRecord).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &apiKeyRecord, nil
}

func (r *tenantWebhookApiKeyRepository) GetFirstApiKeyForTenant(ctx context.Context, tenant string) (*entity.TenantWebhookApiKey, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWebhookApiKeyRepository.GetFirstApiKeyForTenant")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	// get record for tenant or nil if not found
	var apiKeyRecord entity.TenantWebhookApiKey
	err := r.gormDb.
		Where("tenant_name = ?", tenant).
		Order("created_at ASC").
		First(&apiKeyRecord).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.LogFields(log.Bool("result.found", false))
			return nil, nil
		}
		return nil, err
	}
	span.LogFields(log.Bool("result.found", true))
	return &apiKeyRecord, nil
}
