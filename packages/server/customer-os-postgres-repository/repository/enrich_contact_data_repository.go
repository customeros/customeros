package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository/helper"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type enrichContactDataRepository struct {
	gormDb *gorm.DB
}

type EnrichContactDataRepository interface {
	Add(ctx context.Context, data entity.EnrichContactData) helper.QueryResult
}

func NewEnrichContactDataRepository(gormDb *gorm.DB) EnrichContactDataRepository {
	return &enrichContactDataRepository{gormDb: gormDb}
}

func (e enrichContactDataRepository) Add(ctx context.Context, data entity.EnrichContactData) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichContactDataRepository.Add")
	defer span.Finish()

	err := e.gormDb.Create(&data).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: data}
}
