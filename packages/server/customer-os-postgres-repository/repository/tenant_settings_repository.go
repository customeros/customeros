package repository

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"log"
)

type TenantSettingsRepository interface {
	FindForTenantName(ctx context.Context, tenantName string) (*entity.TenantSettings, error)
	Save(ctx context.Context, tenantSettings *entity.TenantSettings) (*entity.TenantSettings, error)
	CheckKeysExist(ctx context.Context, tenantName string, keyName []string) (bool, error)
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

func (r *tenantSettingsRepo) CheckKeysExist(ctx context.Context, tenantName string, keyName []string) (bool, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantSettingsRepository.CheckKeysExist")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var rows int64
	exists := true
	for _, key := range keyName {
		log.Printf("CheckKeysExist: %s, %s", tenantName, key)
		err := r.db.Model(&entity.GoogleServiceAccountKey{}).
			Where(&entity.GoogleServiceAccountKey{TenantName: tenantName, Key: key}, "tenant_name", "key").Count(&rows).Error

		if err != nil {
			return false, fmt.Errorf("CheckKeysExist: %w", err)
		}
		if rows == 0 {
			exists = false
		}

	}
	return exists, nil
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
