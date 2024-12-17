package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowActionExecutionRepository interface {
	Create(ctx context.Context, executionRecord entity.FlowActionExecution) (*entity.FlowActionExecution, error)
	Find(ctx context.Context, executionRecord entity.FlowActionExecution) (*entity.FlowActionExecution, error)
	Update(ctx context.Context, executionRecord entity.FlowActionExecution) (*entity.FlowActionExecution, error)
}

type flowActionExecutionRepository struct {
	gormDb *gorm.DB
}

func NewFlowActionExecutionRepository(gormDb *gorm.DB) FlowActionExecutionRepository {
	return &flowActionExecutionRepository{gormDb: gormDb}
}

func (f *flowActionExecutionRepository) Create(ctx context.Context, executionRecord entity.FlowActionExecution) (*entity.FlowActionExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if executionRecord.Action == "" || executionRecord.FlowExecutionID == "" {
		span.LogFields(log.Object("executionRecord", executionRecord))
		err := errors.New("Action or FlowExecutionID missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	err := f.gormDb.Create(&executionRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &executionRecord, nil
}

func (f *flowActionExecutionRepository) Find(ctx context.Context, executionRecord entity.FlowActionExecution) (*entity.FlowActionExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var flowActionExecution entity.FlowActionExecution
	err := f.gormDb.
		Where(&executionRecord).
		First(&flowActionExecution).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &flowActionExecution, nil
}

func (f *flowActionExecutionRepository) Update(ctx context.Context, executionRecord entity.FlowActionExecution) (*entity.FlowActionExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if executionRecord.ID == "" {
		err := errors.New("ID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedRecord entity.FlowActionExecution
	err := f.gormDb.Model(&executionRecord).Updates(&executionRecord).First(&updatedRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedRecord, nil
}
