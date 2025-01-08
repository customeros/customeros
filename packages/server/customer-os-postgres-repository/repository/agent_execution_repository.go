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

type AgentExecutionRepository interface {
	Create(ctx context.Context, executionRecord entity.AgentExecution) (*entity.AgentExecution, error)
	Find(ctx context.Context, executionRecord entity.AgentExecution) (*entity.AgentExecution, error)
	Update(ctx context.Context, executionRecord entity.AgentExecution) (*entity.AgentExecution, error)
}

type flowAgentExecutionRepository struct {
	gormDb *gorm.DB
}

func NewAgentExecutionRepository(gormDb *gorm.DB) AgentExecutionRepository {
	return &flowAgentExecutionRepository{gormDb: gormDb}
}

func (f *flowAgentExecutionRepository) Create(ctx context.Context, executionRecord entity.AgentExecution) (*entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if executionRecord.AgentID == "" || executionRecord.AutomationExecutionID == "" {
		span.LogFields(log.Object("executionRecord", executionRecord))
		err := errors.New("Agent or ExecutionID missing")
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

func (f *flowAgentExecutionRepository) Find(ctx context.Context, executionRecord entity.AgentExecution) (*entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var flowAgentExecution entity.AgentExecution
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

func (f *flowAgentExecutionRepository) Update(ctx context.Context, executionRecord entity.AgentExecution) (*entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if executionRecord.ID == "" {
		err := errors.New("ID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedRecord entity.AgentExecution
	err := f.gormDb.Model(&executionRecord).Updates(&executionRecord).First(&updatedRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedRecord, nil
}
