package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository/helper"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TenantWebhookRepo struct {
	db *gorm.DB
}

type TenantWebhookRepository interface {
	GetWebhook(ctx context.Context, tenant, event string) helper.QueryResult
	GetWebhooks(ctx context.Context, tenant string) helper.QueryResult
	CreateWebhook(ctx context.Context, integration entity.TenantWebhook) helper.QueryResult
}

func NewTenantWebhookRepo(db *gorm.DB) *TenantWebhookRepo {
	return &TenantWebhookRepo{db: db}
}

func (r *TenantWebhookRepo) GetWebhook(ctx context.Context, tenant, event string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantWebhookRepo.GetWebhook")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	var webhookEntity entity.TenantWebhook
	err := r.db.
		Where("tenant_name = ?", tenant).
		Where("event = ?", event).
		First(&webhookEntity).Error

	// Check if the error is ErrRecordNotFound, treat it as a non-error case
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.QueryResult{Result: nil} // Record not found, return nil without error
		}
		// For all other errors, return the error
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &webhookEntity}
}

func (r *TenantWebhookRepo) GetWebhooks(ctx context.Context, tenant string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantWebhookRepo.GetWebhooks")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	var webhookEntity entity.TenantWebhook
	err := r.db.
		Where("tenant_name = ?", tenant).
		Find(&webhookEntity).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &webhookEntity}
}

func (r *TenantWebhookRepo) CreateWebhook(ctx context.Context, webhook entity.TenantWebhook) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantWebhookRepo.CreateWebhook")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	webhookEntity := entity.TenantWebhook{
		TenantName:      webhook.TenantName,
		ApiKey:          webhook.ApiKey,
		Event:           webhook.Event,
		WebhookUrl:      webhook.WebhookUrl,
		AuthHeaderName:  webhook.AuthHeaderName,
		AuthHeaderValue: webhook.AuthHeaderValue,
		UserId:          webhook.UserId,
		UserFirstName:   webhook.UserFirstName,
		UserLastName:    webhook.UserLastName,
		UserEmail:       webhook.UserEmail,
	}

	err := r.db.Create(&webhookEntity).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &webhookEntity}
}
