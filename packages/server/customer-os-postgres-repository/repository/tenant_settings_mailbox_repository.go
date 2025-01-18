package postgres_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type tenantSettingsMailboxRepository struct {
	gormDb *gorm.DB
}

type TenantSettingsMailboxRepository interface {
	GetAll(ctx context.Context) ([]*postgres_entity.TenantSettingsMailbox, error)
	GetForRampUp(ctx context.Context) ([]*postgres_entity.TenantSettingsMailbox, error)
	GetById(ctx context.Context, id string) (*postgres_entity.TenantSettingsMailbox, error)
	GetByMailbox(ctx context.Context, mailbox string) (*postgres_entity.TenantSettingsMailbox, error)
	GetAllByDomain(ctx context.Context, domain string) ([]postgres_entity.TenantSettingsMailbox, error)
	GetAllByUserId(ctx context.Context, userId string) ([]postgres_entity.TenantSettingsMailbox, error)

	Merge(ctx context.Context, tx *gorm.DB, mailbox *postgres_entity.TenantSettingsMailbox) error
}

func NewTenantSettingsMailboxRepository(db *gorm.DB) TenantSettingsMailboxRepository {
	return &tenantSettingsMailboxRepository{gormDb: db}
}

func (r *tenantSettingsMailboxRepository) GetAll(ctx context.Context) ([]*postgres_entity.TenantSettingsMailbox, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.GetAll")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	var result []*postgres_entity.TenantSettingsMailbox
	err := r.gormDb.
		Where("tenant = ?", tenant).
		Find(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, nil
}

func (r *tenantSettingsMailboxRepository) GetForRampUp(ctx context.Context) ([]*postgres_entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.GetForRampUp")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	var result []*postgres_entity.TenantSettingsMailbox
	err := r.gormDb.
		Where("ramp_up_current < ramp_up_max and (last_ramp_up_at is null or last_ramp_up_at < ?)", utils.StartOfDayInUTC(utils.Now())).
		Find(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, nil
}

func (r *tenantSettingsMailboxRepository) GetById(ctx context.Context, id string) (*postgres_entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogFields(tracingLog.String("id", id))

	var result postgres_entity.TenantSettingsMailbox
	err := r.gormDb.
		Where("tenant = ? and id = ?", tenant, id).
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

func (r *tenantSettingsMailboxRepository) GetByMailbox(ctx context.Context, mailbox string) (*postgres_entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.GetByMailbox")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogFields(tracingLog.String("mailbox", mailbox))

	var result postgres_entity.TenantSettingsMailbox
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

func (r *tenantSettingsMailboxRepository) GetAllByDomain(ctx context.Context, domain string) ([]postgres_entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.GetAllByDomain")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogKV("domain", domain)

	var result []postgres_entity.TenantSettingsMailbox
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

func (r *tenantSettingsMailboxRepository) GetAllByUserId(ctx context.Context, userId string) ([]postgres_entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.GetAllByUserId")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogKV("userId", userId)

	var result []postgres_entity.TenantSettingsMailbox
	err := r.gormDb.WithContext(ctx).
		Where("tenant = ? and user_id = ?", tenant, userId).
		Find(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "db error"))
		return nil, err
	}

	return result, nil
}

func (r *tenantSettingsMailboxRepository) Merge(ctx context.Context, tx *gorm.DB, input *postgres_entity.TenantSettingsMailbox) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogFields(tracingLog.Object("mailbox", input))

	// Check if the mailbox already exists
	var mailbox postgres_entity.TenantSettingsMailbox
	err := r.gormDb.
		Where("tenant = ? AND mailbox_username = ?", tenant, input.MailboxUsername).
		First(&mailbox).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tracing.TraceErr(span, err)
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// If not found, create a new mailbox
		mailbox = postgres_entity.TenantSettingsMailbox{
			Tenant:                  tenant,
			MailboxUsername:         input.MailboxUsername,
			MailboxPassword:         input.MailboxPassword,
			Status:                  input.Status,
			ForwardingTo:            input.ForwardingTo,
			WebmailEnabled:          input.WebmailEnabled,
			Username:                input.Username,
			UserId:                  input.UserId,
			Domain:                  input.Domain,
			LastRampUpAt:            utils.Now(),
			RampUpRate:              3,
			RampUpMax:               40,
			RampUpCurrent:           3,
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
		mailbox.Status = input.Status
		mailbox.ForwardingTo = input.ForwardingTo
		mailbox.WebmailEnabled = input.WebmailEnabled
		mailbox.MailboxPassword = input.MailboxPassword
		mailbox.LastRampUpAt = input.LastRampUpAt
		mailbox.RampUpRate = input.RampUpRate
		mailbox.RampUpMax = input.RampUpMax
		mailbox.RampUpCurrent = input.RampUpCurrent
		mailbox.MinMinutesBetweenEmails = input.MinMinutesBetweenEmails
		mailbox.MaxMinutesBetweenEmails = input.MaxMinutesBetweenEmails
		mailbox.UserId = input.UserId
		mailbox.UpdatedAt = utils.Now()

		if tx == nil {
			tx = r.gormDb
		}

		err = tx.Save(&mailbox).Error
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
