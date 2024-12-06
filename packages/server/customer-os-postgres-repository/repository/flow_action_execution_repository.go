package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type FlowActionExecutionRepository interface {
	Save(ctx context.Context, executionRecord entity.ActionExecution) (string, error)
}

type flowActionExecutionRepository struct {
	gormDb *gorm.DB
}

func NewFlowActionExecutionRepository(gormDb *gorm.DB) FlowActionExecutionRepository {
	return &flowActionExecutionRepository{gormDb: gormDb}
}

func (f *flowActionExecutionRepository) Save(ctx context.Context, executionRecord entity.ActionExecution) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if executionRecord.Action == "" || executionRecord.FlowExecutionID == "" {
		span.LogFields(log.Object("executionRecord", executionRecord))
		err := errors.New("Action or FlowExectionID missing")
		tracing.TraceErr(span, err)
		return "", err
	}

	err := f.gormDb.Save(&executionRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return executionRecord.ID, nil
}
