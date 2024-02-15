package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	"gorm.io/gorm"
)

type TenantWebhookRepo struct {
	db *gorm.DB
}

type TenantWebhookRepository interface {
	GetWebhook(tenant, webhookUrl, event string) helper.QueryResult
	GetWebhooks(tenant string) helper.QueryResult
	SaveWebhook(integration entity.TenantWebhook) helper.QueryResult
}

func NewTenantWebhookRepo(db *gorm.DB) *TenantWebhookRepo {
	return &TenantWebhookRepo{db: db}
}

func (r *TenantWebhookRepo) GetWebhook(tenant, webhookUrl, event string) helper.QueryResult {
	var webhookEntity entity.TenantWebhook
	err := r.db.
		Where("tenant_name = ?", tenant).
		Where("webhook_url = ?", webhookUrl).
		Where("event = ?", event).
		First(&webhookEntity).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &webhookEntity}
}

func (r *TenantWebhookRepo) GetWebhooks(tenant string) helper.QueryResult {
	var webhookEntity entity.TenantWebhook
	err := r.db.
		Where("tenant_name = ?", tenant).
		Find(&webhookEntity).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &webhookEntity}
}

func (r *TenantWebhookRepo) SaveWebhook(webhook entity.TenantWebhook) helper.QueryResult {
	webhookEntity := entity.TenantWebhook{
		TenantName: webhook.TenantName,
		ApiKey:     webhook.ApiKey,
		Event:      webhook.Event,
		WebhookUrl: webhook.WebhookUrl,
	}

	err := r.db.Create(&webhookEntity).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &webhookEntity}
}
