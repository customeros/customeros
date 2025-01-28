package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type QuickbooksSettingsRepository interface {
	Get(ctx context.Context, tenant string) (*postgres_entity.QuickbooksSettingsEntity, error)
	Save(ctx context.Context, quickbooksSettings postgres_entity.QuickbooksSettingsEntity) (*postgres_entity.QuickbooksSettingsEntity, error)
	Delete(ctx context.Context, tenant string) error
}

type quickbooksSettingsRepository struct {
	db *gorm.DB
}

func NewQuickbooksSettingsRepository(db *gorm.DB) QuickbooksSettingsRepository {
	return &quickbooksSettingsRepository{
		db: db,
	}
}

func (repo *quickbooksSettingsRepository) Get(ctx context.Context, tenant string) (*postgres_entity.QuickbooksSettingsEntity, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "QuickbooksSettingsRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	var existing *postgres_entity.QuickbooksSettingsEntity
	err := repo.db.Find(&existing, "tenant = ?", tenant).Error

	if err != nil {
		span.LogFields(tracingLog.Bool("result.found", false))
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	if existing == nil {
		span.LogFields(tracingLog.Bool("result.found", false))
		return nil, nil
	} else if existing.Tenant == "" {
		span.LogFields(tracingLog.Bool("result.found", false))
		return nil, nil
	}
	span.LogFields(tracingLog.Bool("result.found", true))
	return existing, nil
}

func (repo *quickbooksSettingsRepository) Save(ctx context.Context, quickbooksSettings postgres_entity.QuickbooksSettingsEntity) (*postgres_entity.QuickbooksSettingsEntity, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "QuickbooksSettingsRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	result := repo.db.Save(&quickbooksSettings)
	if result.Error != nil {
		tracing.TraceErr(span, result.Error)
		return nil, result.Error
	}
	return &quickbooksSettings, nil
}

func (repo *quickbooksSettingsRepository) Delete(ctx context.Context, tenant string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "QuickbooksSettingsRepository.Delete")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	existing, err := repo.Get(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = repo.db.Delete(&existing).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
