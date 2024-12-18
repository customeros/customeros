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

type FlowAgentExecutionRepository interface {
	Create(ctx context.Context, executionRecord entity.FlowAgentExecution) (*entity.FlowAgentExecution, error)
	Find(ctx context.Context, executionRecord entity.FlowAgentExecution) (*entity.FlowAgentExecution, error)
	Update(ctx context.Context, executionRecord entity.FlowAgentExecution) (*entity.FlowAgentExecution, error)
}

type flowAgentExecutionRepository struct {
	gormDb *gorm.DB
}

func NewFlowAgentExecutionRepository(gormDb *gorm.DB) FlowAgentExecutionRepository {
	return &flowAgentExecutionRepository{gormDb: gormDb}
}

func (f *flowAgentExecutionRepository) Create(ctx context.Context, executionRecord entity.FlowAgentExecution) (*entity.FlowAgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowAgentExecutionRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if executionRecord.Agent == "" || executionRecord.FlowExecutionID == "" {
		span.LogFields(log.Object("executionRecord", executionRecord))
		err := errors.New("Agent or FlowExecutionID missing")
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

func (f *flowAgentExecutionRepository) Find(ctx context.Context, executionRecord entity.FlowAgentExecution) (*entity.FlowAgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowAgentExecutionRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var flowAgentExecution entity.FlowAgentExecution
	err := f.gormDb.
		Where(&executionRecord).
		First(&flowAgentExecution).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &flowAgentExecution, nil
}

func (f *flowAgentExecutionRepository) Update(ctx context.Context, executionRecord entity.FlowAgentExecution) (*entity.FlowAgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowAgentExecutionRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if executionRecord.ID == "" {
		err := errors.New("ID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedRecord entity.FlowAgentExecution
	err := f.gormDb.Model(&executionRecord).Updates(&executionRecord).First(&updatedRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedRecord, nil
}
