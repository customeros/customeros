package postgres_repository

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type AgentExecutionRepository interface {
	Create(ctx context.Context, executionRecord postgres_entity.AgentExecution) (*postgres_entity.AgentExecution, error)
	Find(ctx context.Context, executionRecord postgres_entity.AgentExecution) (*postgres_entity.AgentExecution, error)
	Update(ctx context.Context, executionID string, completedAt *time.Time, errorMessage string, goalAchieved bool) (*postgres_entity.AgentExecution, error)
	GetById(ctx context.Context, executionID string) (*postgres_entity.AgentExecution, error)
}

type agentExecutionRepository struct {
	gormDb *gorm.DB
}

func NewAgentExecutionRepository(gormDb *gorm.DB) AgentExecutionRepository {
	return &agentExecutionRepository{gormDb: gormDb}
}

func (f *agentExecutionRepository) Create(ctx context.Context, executionRecord postgres_entity.AgentExecution) (*postgres_entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if executionRecord.AgentID == nil {
		span.LogFields(log.Object("executionRecord", executionRecord))
		err := errors.New("Agent or ExecutionID missing")
		tracing.TraceErr(span, err)
		return nil, err
	}
	if executionRecord.Tenant == "" {
		executionRecord.Tenant = common.GetTenantFromContext(ctx)
	}

	err := f.gormDb.Create(&executionRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &executionRecord, nil
}

func (f *agentExecutionRepository) Find(ctx context.Context, executionRecord postgres_entity.AgentExecution) (*postgres_entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var agentExecution postgres_entity.AgentExecution
	err := f.gormDb.
		Where(&executionRecord).
		First(&agentExecution).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &agentExecution, nil
}

func (f *agentExecutionRepository) Update(ctx context.Context, executionID string, completedAt *time.Time, errorMessage string, goalAchieved bool) (*postgres_entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.UpdateToCompleted")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if executionID == "" {
		err := errors.New("ID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	status := enum.AgentExecutionCompleted.String()
	if errorMessage != "" {
		status = enum.AgentExecutionFail.String()
	}

	var updatedRecord postgres_entity.AgentExecution
	err := f.gormDb.
		Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Updates(map[string]interface{}{
			"status":        status,
			"completed_at":  completedAt,
			"error_message": errorMessage,
			"goalAchieved":  goalAchieved,
		}).
		First(&updatedRecord, "id = ?", executionID).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedRecord, nil
}

func (f *agentExecutionRepository) GetById(ctx context.Context, executionID string) (*postgres_entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("executionID", executionID))

	var agentExecution postgres_entity.AgentExecution
	err := f.gormDb.
		Where("id = ?", executionID).
		First(&agentExecution).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &agentExecution, nil
}
