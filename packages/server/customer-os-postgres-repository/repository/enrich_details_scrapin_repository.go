package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type enrichDetailsScrapInRepository struct {
	db *gorm.DB
}

type EnrichDetailsScrapInRepository interface {
	Create(ctx context.Context, data entity.EnrichDetailsScrapIn) (*entity.EnrichDetailsScrapIn, error)
	GetAllByParam1AndFlow(ctx context.Context, param string, flow entity.ScrapInFlow) ([]entity.EnrichDetailsScrapIn, error)
	GetLatestByParam1AndFlow(ctx context.Context, param string, flow entity.ScrapInFlow) (*entity.EnrichDetailsScrapIn, error)
	GetLatestByParam1AndFlowWithPersonFound(ctx context.Context, param string, flow entity.ScrapInFlow) (*entity.EnrichDetailsScrapIn, error)
	GetLatestByParam1AndFlowWithCompanyFound(ctx context.Context, param string, flow entity.ScrapInFlow) (*entity.EnrichDetailsScrapIn, error)
	GetLatestByAllParamsAndFlow(ctx context.Context, param1, param2, param3, param4 string, flow entity.ScrapInFlow) (*entity.EnrichDetailsScrapIn, error)
	GetLatestByAllParamsAndFlowWithPersonFound(ctx context.Context, param1, param2, param3, param4 string, flow entity.ScrapInFlow) (*entity.EnrichDetailsScrapIn, error)
	GetById(ctx context.Context, id uint64) (*entity.EnrichDetailsScrapIn, error)
}

func NewEnrichDetailsScrapInRepository(gormDb *gorm.DB) EnrichDetailsScrapInRepository {
	return &enrichDetailsScrapInRepository{db: gormDb}
}

func (r enrichDetailsScrapInRepository) GetAllByParam1AndFlow(ctx context.Context, param string, flow entity.ScrapInFlow) ([]entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetAllByParam1AndFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data []entity.EnrichDetailsScrapIn
	err := r.db.Where("param1 = ? AND flow = ?", param, flow).Find(&data).Error
	if err != nil {
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(data)))

	return data, nil
}

func (r enrichDetailsScrapInRepository) GetLatestByParam1AndFlow(ctx context.Context, param string, flow entity.ScrapInFlow) (*entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data entity.EnrichDetailsScrapIn
	err := r.db.Where("param1 = ? AND flow = ?", param, flow).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r enrichDetailsScrapInRepository) Create(ctx context.Context, data entity.EnrichDetailsScrapIn) (*entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "data", data)

	data.CreatedAt = utils.Now()
	data.UpdatedAt = utils.Now()
	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r enrichDetailsScrapInRepository) GetById(ctx context.Context, id uint64) (*entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data entity.EnrichDetailsScrapIn
	err := r.db.Where("id = ?", id).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r enrichDetailsScrapInRepository) GetLatestByAllParamsAndFlow(ctx context.Context, param1, param2, param3, param4 string, flow entity.ScrapInFlow) (*entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetLatestByAllParamsAndFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data entity.EnrichDetailsScrapIn
	err := r.db.Where("param1 = ? AND param2 = ? AND param3 = ? AND param4 = ? AND flow = ?", param1, param2, param3, param4, flow).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r enrichDetailsScrapInRepository) GetLatestByParam1AndFlowWithPersonFound(ctx context.Context, param string, flow entity.ScrapInFlow) (*entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetLatestByParam1AndFlowWithPersonFound")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data entity.EnrichDetailsScrapIn
	err := r.db.Where("param1 = ? AND flow = ? AND person_found = ?", param, flow, true).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r enrichDetailsScrapInRepository) GetLatestByParam1AndFlowWithCompanyFound(ctx context.Context, param string, flow entity.ScrapInFlow) (*entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetLatestByParam1AndFlowWithCompanyFound")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data entity.EnrichDetailsScrapIn
	err := r.db.Where("param1 = ? AND flow = ? AND person_found = ?", param, flow, true).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r enrichDetailsScrapInRepository) GetLatestByAllParamsAndFlowWithPersonFound(ctx context.Context, param1, param2, param3, param4 string, flow entity.ScrapInFlow) (*entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetLatestByAllParamsAndFlowWithPersonFound")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data entity.EnrichDetailsScrapIn
	err := r.db.Where("param1 = ? AND param2 = ? AND param3 = ? AND param4 = ? AND flow = ? AND person_found = ?", param1, param2, param3, param4, flow, true).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}
