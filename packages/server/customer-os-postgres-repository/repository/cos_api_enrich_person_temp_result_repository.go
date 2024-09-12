package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type cosApiEnrichPersonTempResultRepository struct {
	db *gorm.DB
}

type CosApiEnrichPersonTempResultRepository interface {
	Create(ctx context.Context, data entity.CosApiEnrichPersonTempResult) (*entity.CosApiEnrichPersonTempResult, error)
	GetById(ctx context.Context, id, tenant string) (*entity.CosApiEnrichPersonTempResult, error)
	GetByBettercontactRecordId(ctx context.Context, bettercontactRecordId string) (*entity.CosApiEnrichPersonTempResult, error)
}

func NewCosApiEnrichPersonTempResultRepository(gormDb *gorm.DB) CosApiEnrichPersonTempResultRepository {
	return &cosApiEnrichPersonTempResultRepository{db: gormDb}
}

func (r cosApiEnrichPersonTempResultRepository) GetById(ctx context.Context, id, tenant string) (*entity.CosApiEnrichPersonTempResult, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CosApiEnrichPersonTempResultRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data entity.CosApiEnrichPersonTempResult
	err := r.db.Where("id = ? AND tenant = ?", id, tenant).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (r cosApiEnrichPersonTempResultRepository) Create(ctx context.Context, data entity.CosApiEnrichPersonTempResult) (*entity.CosApiEnrichPersonTempResult, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CosApiEnrichPersonTempResultRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "data", data)

	data.CreatedAt = utils.Now()
	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r cosApiEnrichPersonTempResultRepository) GetByBettercontactRecordId(ctx context.Context, bettercontactRecordId string) (*entity.CosApiEnrichPersonTempResult, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CosApiEnrichPersonTempResultRepository.GetByBettercontactRecordId")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data entity.CosApiEnrichPersonTempResult
	err := r.db.Where("bettercontact_record_id = ?", bettercontactRecordId).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}
