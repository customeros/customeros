package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type tenantSettingsMailboxRepository struct {
	gormDb *gorm.DB
}

type TenantSettingsMailboxRepository interface {
	Get(c context.Context, tenant string) ([]*entity.TenantSettingsMailbox, error)
	GetByMailbox(c context.Context, tenant, mailbox string) (*entity.TenantSettingsMailbox, error)
	GetById(c context.Context, tenant, id string) (*entity.TenantSettingsMailbox, error)
	SaveMailbox(c context.Context, tenant, mailboxEmail, mailboxPassword string) error
}

func NewTenantSettingsMailboxRepository(db *gorm.DB) TenantSettingsMailboxRepository {
	return &tenantSettingsMailboxRepository{gormDb: db}
}

func (r *tenantSettingsMailboxRepository) Get(c context.Context, tenant string) ([]*entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(c, "TenantSettingsMailboxRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tracing.TagTenant(span, tenant)
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
	tracing.TagComponentPostgresRepository(span)

	tracing.TagTenant(span, tenant)
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
	tracing.TagComponentPostgresRepository(span)

	tracing.TagTenant(span, tenant)
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

func (r *tenantSettingsMailboxRepository) SaveMailbox(c context.Context, tenant, mailboxEmail, mailboxPassword string) error {
	span, _ := opentracing.StartSpanFromContext(c, "TenantSettingsMailboxRepository.SaveMailbox")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogKV("mailboxEmail", mailboxEmail)

	// Check if the mailbox already exists
	var mailbox entity.TenantSettingsMailbox
	err := r.gormDb.
		Where("tenant = ? AND mailbox_username = ?", tenant, mailboxEmail).
		First(&mailbox).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tracing.TraceErr(span, err)
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// If not found, create a new mailbox
		mailbox = entity.TenantSettingsMailbox{
			Tenant:          tenant,
			MailboxUsername: mailboxEmail,
			MailboxPassword: mailboxPassword,
		}

		err = r.gormDb.Create(&mailbox).Error
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		// If found, update the existing mailbox
		mailbox.MailboxPassword = mailboxPassword
		mailbox.UpdatedAt = utils.Now()

		err = r.gormDb.Save(&mailbox).Error
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
