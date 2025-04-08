package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type BrowserConfigRepository interface {
	Get(ctx context.Context) ([]postgres_entity.BrowserConfig, error)
	GetForUser(ctx context.Context, userId string) (*postgres_entity.BrowserConfig, error)

	Merge(ctx context.Context, browserConfig *postgres_entity.BrowserConfig) error
}

type browserConfigRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewBrowserConfigRepository(gormDb *gorm.DB) BrowserConfigRepository {
	return &browserConfigRepositoryImpl{gormDb: gormDb}
}

func (repo *browserConfigRepositoryImpl) Get(ctx context.Context) ([]postgres_entity.BrowserConfig, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "BrowserConfigRepository.Get")
	defer spans.Finish()

	var result []postgres_entity.BrowserConfig
	err := repo.gormDb.Where("session_status = 'VALID'").Find(&result).Error

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return result, nil
}

func (repo *browserConfigRepositoryImpl) GetForUser(ctx context.Context, userId string) (*postgres_entity.BrowserConfig, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "BrowserConfigRepository.Get")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	var result *postgres_entity.BrowserConfig
	err := repo.gormDb.Where("tenant = ? and user_id = ? and session_status = 'VALID'", tenant, userId).First(&result).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return result, nil
}

func (r *browserConfigRepositoryImpl) Merge(ctx context.Context, input *postgres_entity.BrowserConfig) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "BrowserConfigRepository.Merge")
	defer spans.Finish()
	spans.LogObjectAsJson("input", input)

	tenant := common.GetTenantFromContext(ctx)

	// Check if the browserConfig already exists
	var browserConfig postgres_entity.BrowserConfig
	err := r.gormDb.
		Where("tenant = ? AND user_id = ?", tenant, input.UserId).
		First(&browserConfig).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		spans.TraceError(err)
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		browserConfig = postgres_entity.BrowserConfig{
			Tenant: tenant,
			UserId: input.UserId,
			Status: input.Status,
		}

		err = r.gormDb.Create(&browserConfig).Error
		if err != nil {
			spans.TraceError(err)
			return err
		}
	} else {
		browserConfig.UserId = input.UserId
		browserConfig.Status = input.Status

		err = r.gormDb.Save(&browserConfig).Error
		if err != nil {
			spans.TraceError(err)
			return err
		}
	}

	return nil
}
