package postgres_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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
	GetByMailbox(ctx context.Context, mailbox string) (*postgres_entity.TenantSettingsMailbox, error)
	GetAllByUserId(ctx context.Context, userId string) ([]*postgres_entity.TenantSettingsMailbox, error)
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

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		span.LogFields(tracingLog.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(tracingLog.Bool("result.found", true))

	return &result, nil
}

func (r *tenantSettingsMailboxRepository) GetAllByUserId(ctx context.Context, userId string) ([]*postgres_entity.TenantSettingsMailbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsMailboxRepository.GetAllByUserId")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogKV("userId", userId)

	var result []*postgres_entity.TenantSettingsMailbox
	err := r.gormDb.WithContext(ctx).
		Where("tenant = ? and user_id = ?", tenant, userId).
		Find(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(result)))
	return result, nil
}
