package postgres_repository

import (
	"context"
	"errors"
	"time"

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
	Fail(ctx context.Context, executionID, errorMessage string) error
	Pending(ctx context.Context, executionID string) error
	Finish(ctx context.Context, executionID string) error
	Completed(ctx context.Context, executionID string, goalAchieved *bool) (*postgres_entity.AgentExecution, error)
	ScheduleRetry(ctx context.Context, executionID string, err error) error
	SaveAsyncState(ctx context.Context, executionID string, currentStep string, stateData map[string]any) error
	CompleteStep(ctx context.Context, executionID string, step string, result map[string]any) error
	GoalAchieved(ctx context.Context, executionID string, goalAchieved bool) error
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

func (f *agentExecutionRepository) Completed(ctx context.Context, executionID string, goalAchieved *bool) (*postgres_entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.Completed")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.TagEntity(span, executionID)
	span.LogFields(log.String("executionID", executionID))
	if goalAchieved != nil {
		span.LogFields(log.Bool("goalAchieved", *goalAchieved))
	}

	if executionID == "" {
		err := errors.New("ID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// load the record
	fields := map[string]interface{}{
		"status":       enum.AgentExecutionCompleted.String(),
		"completed_at": utils.NowPtr(),
	}
	if goalAchieved != nil {
		fields["goal_achieved"] = *goalAchieved
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

func (f *agentExecutionRepository) Finish(ctx context.Context, executionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.Finish")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("executionID", executionID))

	// only running agent executions can be finished
	err := f.gormDb.
		Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Where("status = ?", enum.AgentExecutionRunning.String()).
		Updates(map[string]interface{}{
			"status": enum.AgentExecutionCompleted.String(),
		}).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (f *agentExecutionRepository) Fail(ctx context.Context, executionID, errorMessage string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.Fail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("executionID", executionID))

	// only running agent executions can be finished
	err := f.gormDb.
		Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Where("status <> ?", enum.AgentExecutionCompleted.String()).
		Updates(map[string]interface{}{
			"status":        enum.AgentExecutionError.String(),
			"error_message": errorMessage,
		}).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (f *agentExecutionRepository) Pending(ctx context.Context, executionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.Pending")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("executionID", executionID))

	if executionID == "" {
		err := errors.New("executionID is missing")
		tracing.TraceErr(span, err)
		return err
	}

	err := f.gormDb.
		Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Where("status <> ?", enum.AgentExecutionCompleted.String()).
		Updates(map[string]interface{}{
			"status": enum.AgentExecutionPending.String(),
			// Don't set completed_at since it's pending
			// Don't clear error_message in case we want to preserve previous errors
		}).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
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
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.TagEntity(span, executionID)

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

func (f *agentExecutionRepository) ScheduleRetry(ctx context.Context, executionID string, inputError error) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.ScheduleRetry")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagEntity(span, executionID)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	span.LogFields(log.String("executionID", executionID))
	if inputError == nil {
		span.LogFields(log.String("inputError", "nil"))
	} else {
		span.LogFields(log.String("inputError", inputError.Error()))
	}

	execution := &postgres_entity.AgentExecution{}
	result := f.gormDb.First(execution, "id = ?", executionID)
	if result.Error != nil {
		return result.Error
	}

	// Increment retry count and set next retry time
	execution.RetryCount++
	if execution.RetryCount > execution.MaxRetries {
		execution.Status = enum.AgentExecutionError
		if inputError != nil {
			execution.ErrorMessage = utils.StringPtr(inputError.Error())
		}
		execution.NextRetryAt = nil
	} else {

		backoff := utils.CalculateExponentialBackoffDelay(execution.RetryCount, utils.BackoffConfig{
			InitialDelay: 30 * time.Second,
			MaxDelay:     12 * time.Hour,
			Factor:       2.0,
			Jitter:       0.12,
		})
		nextRetry := time.Now().Add(backoff)
		execution.NextRetryAt = &nextRetry
		if inputError != nil {
			execution.ErrorMessage = utils.StringPtr(inputError.Error())
		}
		execution.Status = enum.AgentExecutionRetrying
	}

	return f.gormDb.Save(execution).Error
}

func (f *agentExecutionRepository) SaveAsyncState(ctx context.Context, executionID string, currentStep string, stateData map[string]any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.SaveAsyncState")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("executionID", executionID))

	return f.gormDb.Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Updates(map[string]interface{}{
			"current_step": currentStep,
			"state_data":   stateData,
			"status":       enum.AgentExecutionPending,
		}).Error
}

func (f *agentExecutionRepository) CompleteStep(ctx context.Context, executionID string, step string, result map[string]any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.CompleteStep")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("executionID", executionID))

	execution := &postgres_entity.AgentExecution{}
	err := f.gormDb.First(execution, "id = ?", executionID).Error
	if err != nil {
		return err
	}

	// Initialize checkpoints if nil
	if execution.Checkpoints == nil {
		execution.Checkpoints = make(map[string]any)
	}

	// Store step completion data
	execution.Checkpoints[step] = result
	execution.CurrentStep = "" // Clear current step
	execution.StateData = nil  // Clear state data since step is complete

	return f.gormDb.Save(execution).Error
}

func (f *agentExecutionRepository) GoalAchieved(ctx context.Context, executionID string, goalAchieved bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.GoalAchieved")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagEntity(span, executionID)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	span.LogFields(log.Bool("goalAchieved", goalAchieved))

	if executionID == "" {
		err := errors.New("ID is missing")
		tracing.TraceErr(span, err)
		return err
	}

	err := f.gormDb.
		Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Updates(
			map[string]interface{}{
				"goal_achieved": goalAchieved,
				"completed_at":  utils.NowPtr(),
			}).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
