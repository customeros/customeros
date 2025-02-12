package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TenantSettingsRepository interface {
	FindForTenantName(ctx context.Context, tenantName string) (*postgres_entity.TenantSettings, error)
	Save(ctx context.Context, tenantSettings *postgres_entity.TenantSettings) (*postgres_entity.TenantSettings, error)
}

type tenantSettingsRepo struct {
	db *gorm.DB
}

func NewTenantSettingsRepository(db *gorm.DB) TenantSettingsRepository {
	return &tenantSettingsRepo{
		db: db,
	}
}

func (r *tenantSettingsRepo) FindForTenantName(ctx context.Context, tenantName string) (*postgres_entity.TenantSettings, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsRepository.FindForTenantName")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var tenantSettings postgres_entity.TenantSettings

	err := r.db.
		Where("tenant_name = ?", tenantName).
		First(&tenantSettings).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		span.LogFields(tracingLog.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(tracingLog.Bool("result.found", true))

	return &tenantSettings, nil
}

func (r *tenantSettingsRepo) Save(ctx context.Context, tenantSettings *postgres_entity.TenantSettings) (*postgres_entity.TenantSettings, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.db.Save(tenantSettings).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return tenantSettings, nil
}
