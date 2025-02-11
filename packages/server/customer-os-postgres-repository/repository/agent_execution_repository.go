package postgres_repository

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type AgentExecutionRepository interface {
	Create(ctx context.Context, executionRecord postgres_entity.AgentExecution) (*postgres_entity.AgentExecution, error)
	Find(ctx context.Context, executionRecord postgres_entity.AgentExecution) (*postgres_entity.AgentExecution, error)
	GetById(ctx context.Context, executionID string) (*postgres_entity.AgentExecution, error)
	Fail(ctx context.Context, executionID, errorMessage string) (*postgres_entity.AgentExecution, error)
	Finish(ctx context.Context, executionID string) (*postgres_entity.AgentExecution, error)
	Completed(ctx context.Context, executionID string, goalAchieved bool) (*postgres_entity.AgentExecution, error)
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

func (f *agentExecutionRepository) Completed(ctx context.Context, executionID string, goalAchieved bool) (*postgres_entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.Completed")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("executionID", executionID))
	span.LogFields(log.Bool("goalAchieved", goalAchieved))

	if executionID == "" {
		err := errors.New("ID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// load the record
	fields := map[string]interface{}{
		"status":        enum.AgentExecutionCompleted.String(),
		"goal_achieved": goalAchieved,
		"completed_at":  utils.NowPtr(),
	}

	var updatedRecord postgres_entity.AgentExecution
	err := f.gormDb.
		Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Updates(fields).
		First(&updatedRecord, "id = ?", executionID).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedRecord, nil
}

func (f *agentExecutionRepository) Finish(ctx context.Context, executionID string) (*postgres_entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.Finish")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("executionID", executionID))

	var updatedRecord postgres_entity.AgentExecution
	// only running agent executions can be finished
	err := f.gormDb.
		Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Where("status = ?", enum.AgentExecutionRunning.String()).
		Updates(map[string]interface{}{
			"status": enum.AgentExecutionFinished.String(),
		}).
		First(&updatedRecord, "id = ?", executionID).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedRecord, nil
}

func (f *agentExecutionRepository) Fail(ctx context.Context, executionID, errorMessage string) (*postgres_entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.Fail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("executionID", executionID))

	var updatedRecord postgres_entity.AgentExecution
	// only running agent executions can be finished
	err := f.gormDb.
		Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Where("status <> ?", enum.AgentExecutionCompleted.String()).
		Updates(map[string]interface{}{
			"status":        enum.AgentExecutionFail.String(),
			"error_message": errorMessage,
			"goal_achieved": false,
		}).
		First(&updatedRecord, "id = ?", executionID).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedRecord, nil
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
		span.LogFields(log.Bool("result.found", false))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", true))
	return &agentExecution, nil
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
		span.LogFields(log.Bool("result.found", false))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Bool("result.found", true))
	return &agentExecution, nil
}
