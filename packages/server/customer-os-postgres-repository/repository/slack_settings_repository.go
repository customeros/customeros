package repository

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type SlackSettingsRepository interface {
	Get(ctx context.Context, tenant string) (*entity.SlackSettingsEntity, error)
	Save(ctx context.Context, slackSettings entity.SlackSettingsEntity) (*entity.SlackSettingsEntity, error)
	Delete(ctx context.Context, tenant string) error
}

type slackSettingsRepository struct {
	db *gorm.DB
}

func NewSlackSettingsRepository(db *gorm.DB) SlackSettingsRepository {
	return &slackSettingsRepository{
		db: db,
	}
}

func (repo *slackSettingsRepository) Get(ctx context.Context, tenant string) (*entity.SlackSettingsEntity, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SlackSettingsRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	var existing *entity.SlackSettingsEntity
	err := repo.db.Find(&existing, "tenant_name = ?", tenant).Error

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
	} else if existing.TenantName == "" {
		span.LogFields(tracingLog.Bool("result.found", false))
		return nil, nil
	}
	span.LogFields(tracingLog.Bool("result.found", true))
	return existing, nil
}

func (repo *slackSettingsRepository) Save(ctx context.Context, slackSettings entity.SlackSettingsEntity) (*entity.SlackSettingsEntity, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SlackSettingsRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	result := repo.db.Save(&slackSettings)
	if result.Error != nil {
		return nil, fmt.Errorf("saving slack settings failed: %w", result.Error)
	}
	return &slackSettings, nil
}

func (repo *slackSettingsRepository) Delete(ctx context.Context, tenant string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SlackSettingsRepository.Delete")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	existing, err := repo.Get(ctx, tenant)
	if err != nil {
		return err
	}

	err = repo.db.Delete(&existing).Error
	if err != nil {
		return err
	}

	return nil
}
