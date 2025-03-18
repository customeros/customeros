package postgres_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

type tenantSettingsMailboxRepository struct {
	gormDb *gorm.DB
}

type TenantSettingsMailboxRepository interface {
	GetAll(ctx context.Context) ([]*postgres_entity.TenantSettingsMailbox, error)
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
