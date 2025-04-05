package postgres_repository

import (
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
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
	spans, ctx := telemetry.StartPostgresSpan(ctx, "QuickbooksSettingsRepository.Get")
	defer spans.Finish()
	spans.LogKV("tenant", tenant)

	var existing *postgres_entity.QuickbooksSettingsEntity
	err := repo.db.Find(&existing, "tenant = ?", tenant).Error

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
	} else if existing.Tenant == "" {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	spans.LogKV("result.found", true)
	return existing, nil
}

func (repo *quickbooksSettingsRepository) Save(ctx context.Context, quickbooksSettings postgres_entity.QuickbooksSettingsEntity) (*postgres_entity.QuickbooksSettingsEntity, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "QuickbooksSettingsRepository.Save")
	defer spans.Finish()
	spans.LogObjectAsJson("quickbooksSettings", quickbooksSettings)

	quickbooksSettings.UpdatedAt = utils.NowPtr()
	result := repo.db.Save(&quickbooksSettings)
	if result.Error != nil {
		spans.TraceError(result.Error)
		return nil, result.Error
	}
	return &quickbooksSettings, nil
}

func (repo *quickbooksSettingsRepository) Delete(ctx context.Context, tenant string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "QuickbooksSettingsRepository.Delete")
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
