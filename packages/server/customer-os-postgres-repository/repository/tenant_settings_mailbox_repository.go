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
	GetAll(ctx context.Context, tenant string) ([]*entity.TenantSettingsMailbox, error)
	GetByMailbox(ctx context.Context, tenant, mailbox string) (*entity.TenantSettingsMailbox, error)
	GetById(ctx context.Context, tenant, id string) (*entity.TenantSettingsMailbox, error)
	GetAllByDomain(ctx context.Context, tenant, domain string) ([]entity.TenantSettingsMailbox, error)
	GetAllByUsername(ctx context.Context, tenant, username string) ([]entity.TenantSettingsMailbox, error)

	Merge(ctx context.Context, tenant string, mailbox *entity.TenantSettingsMailbox) error
}

func NewTenantSettingsMailboxRepository(db *gorm.DB) TenantSettingsMailboxRepository {
	return &tenantSettingsMailboxRepository{gormDb: db}
}

func (r *tenantSettingsMailboxRepository) GetAll(c context.Context, tenant string) ([]*entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(c, "TenantSettingsMailboxRepository.GetAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

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

func (r *tenantSettingsMailboxRepository) GetByMailbox(ctx context.Context, tenant, mailbox string) (*entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.GetByMailbox")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

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

func (r *tenantSettingsMailboxRepository) GetById(ctx context.Context, tenant, id string) (*entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

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

func (r *tenantSettingsMailboxRepository) GetAllByDomain(ctx context.Context, tenant, domain string) ([]entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.GetAllByDomain")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogKV("domain", domain)

	var result []entity.TenantSettingsMailbox
	err := r.gormDb.WithContext(ctx).
		Where("tenant = ? and domain = ?", tenant, domain).
		Find(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "db error"))
		return nil, err
	}

	return result, nil
}

func (r *tenantSettingsMailboxRepository) GetAllByUsername(ctx context.Context, tenant, username string) ([]entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.GetAllByUsername")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogKV("username", username)

	var result []entity.TenantSettingsMailbox
	err := r.gormDb.WithContext(ctx).
		Where("tenant = ? and user_name = ?", tenant, username).
		Find(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "db error"))
		return nil, err
	}

	return result, nil
}

func (r *tenantSettingsMailboxRepository) Merge(ctx context.Context, tenant string, input *entity.TenantSettingsMailbox) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.Merge")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	span.LogFields(tracingLog.Object("mailbox", input))

	// Check if the mailbox already exists
	var mailbox entity.TenantSettingsMailbox
	err := r.gormDb.
		Where("tenant = ? AND mailbox_username = ?", tenant, input.MailboxUsername).
		First(&mailbox).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tracing.TraceErr(span, err)
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// If not found, create a new mailbox
		mailbox = entity.TenantSettingsMailbox{
			Tenant:                  tenant,
			MailboxUsername:         input.MailboxUsername,
			MailboxPassword:         input.MailboxPassword,
			Username:                input.Username,
			Domain:                  input.Domain,
			MinMinutesBetweenEmails: input.MinMinutesBetweenEmails,
			MaxMinutesBetweenEmails: input.MaxMinutesBetweenEmails,
		}

		err = r.gormDb.Create(&mailbox).Error
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		// If found, update the existing mailbox
		mailbox.MailboxPassword = input.MailboxPassword
		mailbox.MinMinutesBetweenEmails = input.MinMinutesBetweenEmails
		mailbox.MaxMinutesBetweenEmails = input.MaxMinutesBetweenEmails
		mailbox.UpdatedAt = utils.Now()

		err = r.gormDb.Save(&mailbox).Error
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
