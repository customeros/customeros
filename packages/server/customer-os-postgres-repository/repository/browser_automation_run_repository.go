package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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

func (r *browserAutomationRunRepositoryImpl) Get(ctx context.Context, automationType, status string) ([]entity.BrowserAutomationsRun, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "BrowserAutomationRunRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var result []entity.BrowserAutomationsRun
	err := r.gormDb.Where("type = ? and status = ? ", automationType, status).Find(&result).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result)))

	return result, nil
}

func (r *browserAutomationRunRepositoryImpl) Add(ctx context.Context, input *entity.BrowserAutomationsRun) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "BrowserAutomationRunRepository.Add")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.Create(input).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *browserAutomationRunRepositoryImpl) MarkAsProcessed(ctx context.Context, id int) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "BrowserAutomationRunRepository.MarkAsProcessed")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.gormDb.Model(&entity.BrowserAutomationsRun{}).Where("id = ?", id).Update("status", "PROCESSED").Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
