package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type tenantSettingsMailboxRepository struct {
	gormDb *gorm.DB
}

type TenantSettingsMailboxRepository interface {
	Get(c context.Context, tenant string) ([]*entity.TenantSettingsMailbox, error)
	GetByMailbox(c context.Context, tenant, mailbox string) (*entity.TenantSettingsMailbox, error)
	GetById(c context.Context, tenant, id string) (*entity.TenantSettingsMailbox, error)
}

func NewTenantSettingsMailboxRepository(db *gorm.DB) TenantSettingsMailboxRepository {
	return &tenantSettingsMailboxRepository{gormDb: db}
}

func (r *tenantSettingsMailboxRepository) Get(c context.Context, tenant string) ([]*entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(c, "TenantSettingsMailboxRepository.Get")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, tracing.SpanTagComponentPostgresRepository)

	var result []*entity.TenantSettingsMailbox
	err := r.gormDb.
		Where("tenant = ?", tenant).
		Find(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, nil
}

func (r *tenantSettingsMailboxRepository) GetByMailbox(c context.Context, tenant, mailbox string) (*entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(c, "TenantSettingsMailboxRepository.GetByMailbox")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, tracing.SpanTagComponentPostgresRepository)

	span.LogFields(tracingLog.String("mailbox", mailbox))

	var result entity.TenantSettingsMailbox
	err := r.gormDb.
		Where("tenant = ? and mailbox_username = ?", tenant, mailbox).
		First(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			span.LogFields(tracingLog.Bool("result.found", false))
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Bool("result.found", true))

	return &result, nil
}

func (r *tenantSettingsMailboxRepository) GetById(c context.Context, tenant, id string) (*entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(c, "TenantSettingsMailboxRepository.GetById")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, tracing.SpanTagComponentPostgresRepository)

	span.LogFields(tracingLog.String("id", id))

	var result entity.TenantSettingsMailbox
	err := r.gormDb.
		Where("id = ?", id).
		First(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			span.LogFields(tracingLog.Bool("result.found", false))
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Bool("result.found", true))

	return &result, nil
}
