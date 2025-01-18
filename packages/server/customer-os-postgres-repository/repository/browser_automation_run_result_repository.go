package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type BrowserAutomationRunResultRepository interface {
	Get(ctx context.Context, runId int) (*postgres_entity.BrowserAutomationsRunResult, error)
}

type browserAutomationRunResultRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewBrowserAutomationRunResultRepository(gormDb *gorm.DB) BrowserAutomationRunResultRepository {
	return &browserAutomationRunResultRepositoryImpl{gormDb: gormDb}
}

func (repo *browserAutomationRunResultRepositoryImpl) Get(ctx context.Context, runId int) (*postgres_entity.BrowserAutomationsRunResult, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "BrowserAutomationRunRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var result *postgres_entity.BrowserAutomationsRunResult
	err := repo.gormDb.Where("run_id = ? ", runId).Find(&result).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, nil
}
