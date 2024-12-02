package repository

import (
	"context"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowWebhooksRepository interface {
	FindAllActiveWebhooks(ctx context.Context, tenantName string) (int, []entity.FlowWebhooks, error)
	FindActiveWebhook(ctx context.Context, tenantName, integration string) (int, *entity.FlowWebhooks, error)
	DisableWebhook(ctx context.Context, webhookPath string) error
	CreateWebhook(ctx context.Context, webhook *entity.FlowWebhooks) error
	FindWebhookByPath(ctx context.Context, tenantName, webhookPath string) (entity.FlowWebhooks, error)
}

type flowWebhooksRepository struct {
	gormDb *gorm.DB
}

func NewFlowWebhooksRepository(gormDb *gorm.DB) FlowWebhooksRepository {
	return &flowWebhooksRepository{gormDb: gormDb}
}

func (r *flowWebhooksRepository) CreateWebhook(ctx context.Context, webhook *entity.FlowWebhooks) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowWebhooksRepository.CreateWebhook")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if err := r.gormDb.Create(webhook).Error; err != nil {
		return err
	}
	return nil
}

func (r *flowWebhooksRepository) FindWebhookByPath(ctx context.Context, tenantName, webhookPath string) (entity.FlowWebhooks, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowWebhooksRepository.FindWebhookByPath")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var webhook entity.FlowWebhooks
	webhookPath = strings.TrimPrefix(webhookPath, "/")
	err := r.gormDb.
		Where("tenant_name = ? AND webhook_path = ?", tenantName, webhookPath).
		First(&webhook).Error

	return webhook, err
}

func (r *flowWebhooksRepository) FindAllActiveWebhooks(ctx context.Context, tenantName string) (int, []entity.FlowWebhooks, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowWebhooksRepository.FindAllActiveWebhooks")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var webhooks []entity.FlowWebhooks
	err := r.gormDb.
		Where("tenant_name = ? AND enabled = true", tenantName).
		Order("created_at DESC").
		Find(&webhooks).Error

	if err != nil && err.Error() != "record not found" {
		return 0, nil, err
	}

	count := len(webhooks)

	return count, webhooks, nil
}

func (r *flowWebhooksRepository) FindActiveWebhook(ctx context.Context, tenantName, integration string) (int, *entity.FlowWebhooks, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowWebhooksRepository.FindActiveWebhook")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var webhook entity.FlowWebhooks

	err := r.gormDb.
		Where("tenant_name = ? AND integration = ? AND enabled = true", tenantName, integration).
		First(&webhook).Error

	if err != nil && err.Error() != "record not found" {
		return 0, nil, err
	}

	return 1, &webhook, nil
}

func (r *flowWebhooksRepository) DisableWebhook(ctx context.Context, webhookPath string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowWebhooksRepository.DisableWebhook")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	return r.gormDb.
		Where("webhook_path = ? AND enabled = true", webhookPath).
		Update("enabled", false).Error
}
