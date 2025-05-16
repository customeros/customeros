package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository/helper"
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
	CreateWebhook(ctx context.Context, integration postgres_entity.TenantWebhook) helper.QueryResult
}

func NewTenantWebhookRepo(db *gorm.DB) *TenantWebhookRepo {
	return &TenantWebhookRepo{db: db}
}

func (r *TenantWebhookRepo) GetWebhook(ctx context.Context, tenant, event string) helper.QueryResult {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TenantWebhookRepo.GetWebhook")
	defer spans.Finish()

	var webhookEntity postgres_entity.TenantWebhook
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
	spans, _ := telemetry.StartPostgresSpan(ctx, "TenantWebhookRepo.GetWebhooks")
	defer spans.Finish()

	var webhookEntity postgres_entity.TenantWebhook
	err := r.db.
		Where("tenant_name = ?", tenant).
		Find(&webhookEntity).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &webhookEntity}
}

func (r *TenantWebhookRepo) CreateWebhook(ctx context.Context, webhook postgres_entity.TenantWebhook) helper.QueryResult {
	spans, _ := telemetry.StartPostgresSpan(ctx, "TenantWebhookRepo.CreateWebhook")
	defer spans.Finish()

	webhookEntity := postgres_entity.TenantWebhook{
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
