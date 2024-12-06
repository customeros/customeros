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

type FlowExecutionRepository interface {
	Save(ctx context.Context, executionRecord entity.FlowExecution) (string, error)
}

type flowExecutionRepository struct {
	gormDb *gorm.DB
}

func NewFlowExecutionRepository(gormDb *gorm.DB) FlowExecutionRepository {
	return &flowExecutionRepository{gormDb: gormDb}
}

func (f *flowExecutionRepository) Save(ctx context.Context, executionRecord entity.FlowExecution) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if executionRecord.Tenant == "" || executionRecord.FlowID == "" || executionRecord.EntityID == "" {
		span.LogFields(log.Object("executionRecord", executionRecord))
		err := errors.New("Params missing")
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
