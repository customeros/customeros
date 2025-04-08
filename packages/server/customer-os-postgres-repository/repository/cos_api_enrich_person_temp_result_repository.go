package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type cosApiEnrichPersonTempResultRepository struct {
	db *gorm.DB
}

type CosApiEnrichPersonTempResultRepository interface {
	Create(ctx context.Context, data postgres_entity.CosApiEnrichPersonTempResult) (*postgres_entity.CosApiEnrichPersonTempResult, error)
	GetById(ctx context.Context, id, tenant string) (*postgres_entity.CosApiEnrichPersonTempResult, error)
	GetByBettercontactRecordId(ctx context.Context, bettercontactRecordId string) (*postgres_entity.CosApiEnrichPersonTempResult, error)
}

func NewCosApiEnrichPersonTempResultRepository(gormDb *gorm.DB) CosApiEnrichPersonTempResultRepository {
	return &cosApiEnrichPersonTempResultRepository{db: gormDb}
}

func (r cosApiEnrichPersonTempResultRepository) GetById(ctx context.Context, id, tenant string) (*postgres_entity.CosApiEnrichPersonTempResult, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CosApiEnrichPersonTempResultRepository.GetById")
	defer spans.Finish()

	var data postgres_entity.CosApiEnrichPersonTempResult
	err := r.db.Where("id = ? AND tenant = ?", id, tenant).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (r cosApiEnrichPersonTempResultRepository) Create(ctx context.Context, data postgres_entity.CosApiEnrichPersonTempResult) (*postgres_entity.CosApiEnrichPersonTempResult, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CosApiEnrichPersonTempResultRepository.Create")
	defer spans.Finish()
	spans.LogObjectAsJson("data", data)

	data.CreatedAt = utils.Now()
	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r cosApiEnrichPersonTempResultRepository) GetByBettercontactRecordId(ctx context.Context, bettercontactRecordId string) (*postgres_entity.CosApiEnrichPersonTempResult, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "CosApiEnrichPersonTempResultRepository.GetByBettercontactRecordId")
	defer spans.Finish()

	var data postgres_entity.CosApiEnrichPersonTempResult
	err := r.db.Where("bettercontact_record_id = ?", bettercontactRecordId).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	return &data, nil
}
