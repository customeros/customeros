package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository/helper"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type enrichDetailsScrapInRepository struct {
	gormDb *gorm.DB
}

type EnrichDetailsScrapInRepository interface {
	Add(ctx context.Context, data entity.EnrichDetailsScrapIn) helper.QueryResult
	GetAllByParam1AndFlow(ctx context.Context, param string, flow entity.ScrapInFlow) ([]entity.EnrichDetailsScrapIn, error)
}

func NewEnrichDetailsScrapInRepository(gormDb *gorm.DB) EnrichDetailsScrapInRepository {
	return &enrichDetailsScrapInRepository{gormDb: gormDb}
}

func (e enrichDetailsScrapInRepository) GetAllByParam1AndFlow(ctx context.Context, param string, flow entity.ScrapInFlow) ([]entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetAllByParam1AndFlow")
	defer span.Finish()

	var data []entity.EnrichDetailsScrapIn
	err := e.gormDb.Where("param1 = ? AND flow = ?", param, flow).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e enrichDetailsScrapInRepository) Add(ctx context.Context, data entity.EnrichDetailsScrapIn) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.Add")
	defer span.Finish()

	err := e.gormDb.Create(&data).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: data}
}