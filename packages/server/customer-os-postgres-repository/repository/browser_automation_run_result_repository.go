package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
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
	spans, _ := telemetry.StartPostgresSpan(ctx, "BrowserAutomationRunRepository.Get")
	defer spans.Finish()

	var result *postgres_entity.BrowserAutomationsRunResult
	err := repo.gormDb.Where("run_id = ? ", runId).Find(&result).Error

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return result, nil
}
