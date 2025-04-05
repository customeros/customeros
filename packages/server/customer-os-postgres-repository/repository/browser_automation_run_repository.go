package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"golang.org/x/net/context"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type BrowserAutomationRunRepository interface {
	Get(ctx context.Context, automationType, status string) ([]postgres_entity.BrowserAutomationsRun, error)
	Add(ctx context.Context, input *postgres_entity.BrowserAutomationsRun) error
	MarkAsProcessed(ctx context.Context, id int) error
}

type browserAutomationRunRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewBrowserAutomationRunRepository(gormDb *gorm.DB) BrowserAutomationRunRepository {
	return &browserAutomationRunRepositoryImpl{gormDb: gormDb}
}

func (r *browserAutomationRunRepositoryImpl) Get(ctx context.Context, automationType, status string) ([]postgres_entity.BrowserAutomationsRun, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "BrowserAutomationRunRepository.Get")
	defer spans.Finish()

	var result []postgres_entity.BrowserAutomationsRun
	err := r.gormDb.Where("type = ? and status = ? ", automationType, status).Find(&result).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(result))

	return result, nil
}

func (r *browserAutomationRunRepositoryImpl) Add(ctx context.Context, input *postgres_entity.BrowserAutomationsRun) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "BrowserAutomationRunRepository.Add")
	defer spans.Finish()

	err := r.gormDb.Create(input).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *browserAutomationRunRepositoryImpl) MarkAsProcessed(ctx context.Context, id int) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "BrowserAutomationRunRepository.MarkAsProcessed")
	defer spans.Finish()

	err := r.gormDb.Model(&postgres_entity.BrowserAutomationsRun{}).Where("id = ?", id).Update("status", "PROCESSED").Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
