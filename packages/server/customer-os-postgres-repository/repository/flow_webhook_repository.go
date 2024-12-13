package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowWebhooksRepository interface {
	FindAll(ctx context.Context) (*[]entity.FlowWebhooks, error)
	Find(ctx context.Context, webhook entity.FlowWebhooks) (*entity.FlowWebhooks, error)
	Update(ctx context.Context, webhook entity.FlowWebhooks) (*entity.FlowWebhooks, error)
	Create(ctx context.Context, webhook entity.FlowWebhooks) (*entity.FlowWebhooks, error)
}

type flowWebhooksRepository struct {
	gormDb *gorm.DB
}

func NewFlowWebhooksRepository(gormDb *gorm.DB) FlowWebhooksRepository {
	return &flowWebhooksRepository{gormDb: gormDb}
}

func (r *flowWebhooksRepository) Create(ctx context.Context, webhook entity.FlowWebhooks) (*entity.FlowWebhooks, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowWebhooksRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.Create(&webhook).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &webhook, nil
}

func (r *flowWebhooksRepository) FindAll(ctx context.Context) (*[]entity.FlowWebhooks, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowWebhooksRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set in context")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var webhooks []entity.FlowWebhooks
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

func (r *flowWebhooksRepository) Find(ctx context.Context, webhook entity.FlowWebhooks) (*entity.FlowWebhooks, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowWebhooksRepository.Find")
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

	var foundWebhook entity.FlowWebhooks
	query := r.gormDb.Where("enabled = ?", true)

	// Add additional filters based on non-zero fields in webhook
	if webhook != (entity.FlowWebhooks{}) {
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

func (r *flowWebhooksRepository) Update(ctx context.Context, webhook entity.FlowWebhooks) (*entity.FlowWebhooks, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowWebhooksRepository.Update")
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

	var updatedWebhook entity.FlowWebhooks
	err := r.gormDb.Model(&webhook).Updates(&webhook).First(&updatedWebhook).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedWebhook, nil
}
