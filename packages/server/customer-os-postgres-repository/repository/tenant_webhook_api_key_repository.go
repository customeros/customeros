package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TenantWebhookApiKeyRepository interface {
	CreateApiKey(ctx context.Context, tenant string) error
	GetTenantForApiKey(ctx context.Context, apiKey string) (*entity.TenantWebhookApiKey, error)
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
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	apiKey := entity.TenantWebhookApiKey{
		Tenant: tenant,
		Key:    uuid.New().String(),
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
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

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
