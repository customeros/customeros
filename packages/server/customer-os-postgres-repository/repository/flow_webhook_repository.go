package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowWebhooksRepository interface {
	FindAllActiveWebhooks(ctx context.Context, tenantName string) (int64, []entity.FlowWebhooks, error)
	FindActiveWebhook(ctx context.Context, tenantName, integration string) (int64, *entity.FlowWebhooks, error)
	DisableWebhook(ctx context.Context, webhookPath string) error
	CreateWebhook(webhook *entity.FlowWebhooks) error
}

type flowWebhooksRepository struct {
	gormDb *gorm.DB
}

func NewFlowWebhooksRepository(gormDb *gorm.DB) FlowWebhooksRepository {
	return &flowWebhooksRepository{gormDb: gormDb}
}

func (r *flowWebhooksRepository) CreateWebhook(webhook *entity.FlowWebhooks) error {
	if err := r.gormDb.Create(webhook).Error; err != nil {
		return err
	}

	return nil
}

func (r *flowWebhooksRepository) FindAllActiveWebhooks(ctx context.Context, tenantName string) (int64, []entity.FlowWebhooks, error) {
	var count int64

	err := r.gormDb.Model(&entity.FlowWebhooks{}).
		Where("tenant_name = ? AND enabled = true", tenantName).
		Count(&count).Error
	if err != nil {
		return 0, nil, err // DB error
	}

	if count == 0 {
		return 0, nil, nil // No results found
	}

	var webhooks []entity.FlowWebhooks
	err = r.gormDb.
		Where("tenant_name = ? AND enabled = true", tenantName).
		Order("created_at DESC").
		Find(&webhooks).Error
	if err != nil {
		return count, nil, err
	}

	return count, webhooks, nil
}

func (r *flowWebhooksRepository) FindActiveWebhook(ctx context.Context, tenantName, integration string) (int64, *entity.FlowWebhooks, error) {
	var count int64
	var webhook entity.FlowWebhooks

	err := r.gormDb.Model(&entity.FlowWebhooks{}).
		Where("tenant_name = ? AND integration = ? AND enabled = true", tenantName, integration).
		Count(&count).Error
	if err != nil {
		return 0, nil, err // DB error
	}

	if count == 0 {
		return 0, nil, nil // No results found
	}

	err = r.gormDb.
		Where("tenant_name = ? AND integration = ? AND enabled = true", tenantName, integration).
		First(&webhook).Error
	if err != nil {
		return count, nil, err
	}

	return 1, &webhook, nil
}

func (r *flowWebhooksRepository) DisableWebhook(ctx context.Context, webhookPath string) error {
	return r.gormDb.
		Where("webhook_path = ? AND enabled = true", webhookPath).
		Update("enabled", false).Error
}
