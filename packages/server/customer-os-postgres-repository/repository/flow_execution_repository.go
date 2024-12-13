package repository

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

type FlowExecutionRepository interface {
	Create(ctx context.Context, executionRecord entity.FlowExecution) (string, error)
	FindRecord(ctx context.Context, executionRecord entity.FlowExecution) (*entity.FlowExecution, error)
	Update(ctx context.Context, executionRecord entity.FlowExecution) (*entity.FlowExecution, error)
}

type flowExecutionRepository struct {
	gormDb *gorm.DB
}

func NewFlowExecutionRepository(gormDb *gorm.DB) FlowExecutionRepository {
	return &flowExecutionRepository{gormDb: gormDb}
}

func (f *flowExecutionRepository) Create(ctx context.Context, executionRecord entity.FlowExecution) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := f.gormDb.Create(&executionRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return executionRecord.ID, nil
}

func (f *flowExecutionRepository) FindRecord(ctx context.Context, executionRecord entity.FlowExecution) (*entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionRepository.FindRecord")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var foundRecord entity.FlowExecution

	err := f.gormDb.
		Where(&executionRecord).
		First(&foundRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &foundRecord, nil
}

func (f *flowExecutionRepository) Update(ctx context.Context, executionRecord entity.FlowExecution) (*entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var updatedRecord entity.FlowExecution

	err := f.gormDb.Model(&entity.FlowExecution{}).
		Where("id = ?", executionRecord.ID).
		Updates(executionRecord).
		First(&updatedRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedRecord, nil
}
