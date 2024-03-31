package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TenantWebhookRepo struct {
	db *gorm.DB
}

type TenantWebhookRepository interface {
	GetWebhook(tenant, event string) helper.QueryResult
	GetWebhooks(tenant string) helper.QueryResult
	CreateWebhook(integration entity.TenantWebhook) helper.QueryResult
}

func NewTenantWebhookRepo(db *gorm.DB) *TenantWebhookRepo {
	return &TenantWebhookRepo{db: db}
}

func (r *TenantWebhookRepo) GetWebhook(tenant, event string) helper.QueryResult {
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

func (r *TenantWebhookRepo) CreateWebhook(webhook entity.TenantWebhook) helper.QueryResult {
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
