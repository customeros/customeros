package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type WebhooksRepository interface {
	FindAll(ctx context.Context) (*[]postgres_entity.Webhooks, error)
	Find(ctx context.Context, webhook postgres_entity.Webhooks) (*postgres_entity.Webhooks, error)
	Update(ctx context.Context, webhook postgres_entity.Webhooks) (*postgres_entity.Webhooks, error)
	Create(ctx context.Context, webhook postgres_entity.Webhooks) (*postgres_entity.Webhooks, error)
}

type webhooksRepository struct {
	gormDb *gorm.DB
}

func NewWebhooksRepository(gormDb *gorm.DB) WebhooksRepository {
	return &webhooksRepository{gormDb: gormDb}
}

func (r *webhooksRepository) Create(ctx context.Context, webhook postgres_entity.Webhooks) (*postgres_entity.Webhooks, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhooksRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var created postgres_entity.Webhooks
	err := r.gormDb.Create(&webhook).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &webhook, nil
}

func (r *webhooksRepository) FindAll(ctx context.Context) (*[]postgres_entity.Webhooks, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhooksRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set in context")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var webhooks []postgres_entity.Webhooks
	err := r.gormDb.
		Where("enabled = ? AND tenant = ?", true, tenant).
		Order("created_at DESC").
		Find(&webhooks).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &webhooks, nil
}

func (r *webhooksRepository) Find(ctx context.Context, webhook postgres_entity.Webhooks) (*postgres_entity.Webhooks, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhooksRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if webhook.Tenant == "" {
		webhook.Tenant = common.GetTenantFromContext(ctx)
		if webhook.Tenant == "" {
			err := errors.New("Tenant not set in context")
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	var foundWebhook postgres_entity.Webhooks
	query := r.gormDb.Where("enabled = ?", true)

	// Add additional filters based on non-zero fields in webhook
	if webhook != (postgres_entity.Webhooks{}) {
		query = query.Where(&webhook)
	}

	err := query.First(&foundWebhook).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &foundWebhook, nil
}

func (r *webhooksRepository) Update(ctx context.Context, webhook postgres_entity.Webhooks) (*postgres_entity.Webhooks, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebhooksRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if webhook.Tenant == "" {
		webhook.Tenant = common.GetTenantFromContext(ctx)
		if webhook.Tenant == "" {
			err := errors.New("Tenant not set in context")
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	if webhook.WebhookPath == "" {
		err := errors.New("Webhook path is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedWebhook postgres_entity.Webhooks
	err := r.gormDb.Model(&webhook).Updates(&webhook).First(&updatedWebhook).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedWebhook, nil
}
