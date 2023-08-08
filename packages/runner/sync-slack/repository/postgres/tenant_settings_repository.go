package postgres

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/repository/helper"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type TenantSettingsRepository interface {
	FindForTenantName(ctx context.Context, tenantName string) helper.QueryResult
}

type tenantSettingsRepo struct {
	db *gorm.DB
}

func NewTenantSettingsRepository(db *gorm.DB) TenantSettingsRepository {
	return &tenantSettingsRepo{
		db: db,
	}
}

func (r *tenantSettingsRepo) FindForTenantName(ctx context.Context, tenantName string) helper.QueryResult {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsRepo.FindForTenantName")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(span)
	span.LogFields(log.String("tenantName", tenantName))

	var tenantSettings entity.TenantSettings

	err := r.db.
		Where("tenant_name = ?", tenantName).
		First(&tenantSettings).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return helper.QueryResult{Error: err}
	}
	if err == gorm.ErrRecordNotFound {
		return helper.QueryResult{Result: nil}
	}

	return helper.QueryResult{Result: tenantSettings}
}
