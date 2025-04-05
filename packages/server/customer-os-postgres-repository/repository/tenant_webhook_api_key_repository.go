package postgres_repository

import (
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TenantWebhookApiKeyRepository interface {
	CreateApiKey(ctx context.Context, tenant string) error
	GetTenantForApiKey(ctx context.Context, apiKey string) (*postgres_entity.TenantWebhookApiKey, error)
	GetFirstApiKeyForTenant(ctx context.Context, tenant string) (*postgres_entity.TenantWebhookApiKey, error)
}

type tenantWebhookApiKeyRepository struct {
	gormDb *gorm.DB
}

func NewTenantWebhookApiKeyRepository(gormDb *gorm.DB) TenantWebhookApiKeyRepository {
	return &tenantWebhookApiKeyRepository{gormDb: gormDb}
}

func (r *tenantWebhookApiKeyRepository) CreateApiKey(ctx context.Context, tenant string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "TenantWebhookApiKeyRepository.CreateApiKey")
	defer spans.Finish()

	now := utils.Now()
	apiKey := postgres_entity.TenantWebhookApiKey{
		Tenant:    tenant,
		Key:       postgres_entity.KeyPrefix + utils.GenerateKey(32, false),
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := r.gormDb.Create(&apiKey).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}

func (r *tenantWebhookApiKeyRepository) GetTenantForApiKey(ctx context.Context, apiKey string) (*postgres_entity.TenantWebhookApiKey, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "TenantWebhookApiKeyRepository.GetTenantWithApiKey")
	defer spans.Finish()

	// get record for api key or nil if not found
	var apiKeyRecord postgres_entity.TenantWebhookApiKey
	err := r.gormDb.
		Where("key = ?", apiKey).
		First(&apiKeyRecord).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}
	return &apiKeyRecord, nil
}

func (r *tenantWebhookApiKeyRepository) GetFirstApiKeyForTenant(ctx context.Context, tenant string) (*postgres_entity.TenantWebhookApiKey, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "TenantWebhookApiKeyRepository.GetFirstApiKeyForTenant")
	defer spans.Finish()

	// get record for tenant or nil if not found
	var apiKeyRecord postgres_entity.TenantWebhookApiKey
	err := r.gormDb.
		Where("tenant_name = ?", tenant).
		Order("created_at ASC").
		First(&apiKeyRecord).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.found", true)
	return &apiKeyRecord, nil
}
