package repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type TenantSettingsRepository interface {
	FindForTenantName(ctx context.Context, tenantName string) (*entity.TenantSettings, error)
	Save(ctx context.Context, tenantSettings *entity.TenantSettings) (*entity.TenantSettings, error)
}

type tenantSettingsRepo struct {
	db *gorm.DB
}

func NewTenantSettingsRepository(db *gorm.DB) TenantSettingsRepository {
	return &tenantSettingsRepo{
		db: db,
	}
}

func (r *tenantSettingsRepo) FindForTenantName(ctx context.Context, tenantName string) (*entity.TenantSettings, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsRepository.FindForTenantName")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var tenantSettings entity.TenantSettings

	err := r.db.
		Where("tenant_name = ?", tenantName).
		First(&tenantSettings).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &tenantSettings, nil
}

func (r *tenantSettingsRepo) Save(ctx context.Context, tenantSettings *entity.TenantSettings) (*entity.TenantSettings, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.db.Save(tenantSettings).Error

	if err != nil {
		return nil, err
	}

	return tenantSettings, nil
}
