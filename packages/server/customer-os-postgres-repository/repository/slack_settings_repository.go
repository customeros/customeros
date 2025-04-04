package postgres_repository

import (
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type SlackSettingsRepository interface {
	Get(ctx context.Context, tenant string) (*postgres_entity.SlackSettingsEntity, error)
	Save(ctx context.Context, slackSettings postgres_entity.SlackSettingsEntity) (*postgres_entity.SlackSettingsEntity, error)
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

func (repo *slackSettingsRepository) Get(ctx context.Context, tenant string) (*postgres_entity.SlackSettingsEntity, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "SlackSettingsRepository.Get")
	defer spans.Finish()
	spans.LogKV("tenant", tenant)
	var existing *postgres_entity.SlackSettingsEntity
	err := repo.db.Find(&existing, "tenant_name = ?", tenant).Error

	if err != nil {
		spans.LogKV("result.found", false)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	if existing == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	} else if existing.TenantName == "" {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	spans.LogKV("result.found", true)
	return existing, nil
}

func (repo *slackSettingsRepository) Save(ctx context.Context, slackSettings postgres_entity.SlackSettingsEntity) (*postgres_entity.SlackSettingsEntity, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "SlackSettingsRepository.Save")
	defer spans.Finish()

	result := repo.db.Save(&slackSettings)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	return &slackSettings, nil
}

func (repo *slackSettingsRepository) Delete(ctx context.Context, tenant string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "SlackSettingsRepository.Delete")
	defer spans.Finish()
	spans.LogKV("tenant", tenant)

	existing, err := repo.Get(ctx, tenant)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	err = repo.db.Delete(&existing).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
