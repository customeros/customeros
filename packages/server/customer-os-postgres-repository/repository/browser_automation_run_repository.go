package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type BrowserAutomationRunRepository interface {
	Get(ctx context.Context, automationType, status string) ([]entity.BrowserAutomationsRun, error)
	Add(ctx context.Context, input *entity.BrowserAutomationsRun) error
	MarkAsProcessed(ctx context.Context, id int) error
}

type browserAutomationRunRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewBrowserAutomationRunRepository(gormDb *gorm.DB) BrowserAutomationRunRepository {
	return &browserAutomationRunRepositoryImpl{gormDb: gormDb}
}

func (repo *browserAutomationRunRepositoryImpl) Get(ctx context.Context, automationType, status string) ([]entity.BrowserAutomationsRun, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "BrowserAutomationRunRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var result []entity.BrowserAutomationsRun
	err := repo.gormDb.Where("type = ? and status = ? ", automationType, status).Find(&result).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, nil
}

func (repo browserAutomationRunRepositoryImpl) Add(ctx context.Context, input *entity.BrowserAutomationsRun) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "BrowserAutomationRunRepository.Add")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := repo.gormDb.Create(input).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (repo browserAutomationRunRepositoryImpl) MarkAsProcessed(ctx context.Context, id int) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "BrowserAutomationRunRepository.MarkAsProcessed")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := repo.gormDb.Model(&entity.BrowserAutomationsRun{}).Where("id = ?", id).Update("status", "PROCESSED").Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
