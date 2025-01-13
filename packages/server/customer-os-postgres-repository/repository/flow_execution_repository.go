package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

type FlowExecutionRepository interface {
	Create(ctx context.Context, executionRecord entity.FlowExecution) (*entity.FlowExecution, error)
	Find(ctx context.Context, executionRecord entity.FlowExecution) (*entity.FlowExecution, error)
	Update(ctx context.Context, executionRecord entity.FlowExecution) (*entity.FlowExecution, error)
}

type flowExecutionRepository struct {
	gormDb *gorm.DB
}

func NewFlowExecutionRepository(gormDb *gorm.DB) FlowExecutionRepository {
	return &flowExecutionRepository{gormDb: gormDb}
}

func (f *flowExecutionRepository) Create(ctx context.Context, executionRecord entity.FlowExecution) (*entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := f.gormDb.Create(&executionRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &executionRecord, nil
}

func (f *flowExecutionRepository) Find(ctx context.Context, executionRecord entity.FlowExecution) (*entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var foundRecord entity.FlowExecution
	err := f.gormDb.
		Where(&executionRecord).
		First(&foundRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &foundRecord, nil
}

func (f *flowExecutionRepository) FindAll(ctx context.Context, executionRecord entity.FlowExecution) (*[]entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var records []entity.FlowExecution
	err := f.gormDb.
		Where(&executionRecord).
		Find(&records).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &records, nil
}

func (f *flowExecutionRepository) Update(ctx context.Context, executionRecord entity.FlowExecution) (*entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if executionRecord.ID == "" {
		err := errors.New("flow execution ID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedRecord entity.FlowExecution
	err := f.gormDb.
		Model(&entity.FlowExecution{}).
		Where("id = ?", executionRecord.ID).
		Updates(executionRecord).
		First(&updatedRecord, "id = ?", executionRecord.ID).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedRecord, nil
}
