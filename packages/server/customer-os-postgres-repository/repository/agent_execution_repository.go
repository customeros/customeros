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
	ScheduleRetry(ctx context.Context, executionID string, err error, stateData map[string]any) (*time.Time, error)
	SaveAsyncState(ctx context.Context, executionID string, currentStep string, stateData map[string]any) error
	CompleteStep(ctx context.Context, executionID string, step string, result map[string]any) error
	GoalAchieved(ctx context.Context, executionID string, goalAchieved bool, impactedId *string) error
	GetGoalAchievedCountLast30Days(ctx context.Context, agentID string) (int64, error)
	GetExecutionsForRetry(ctx context.Context, limit int) ([]postgres_entity.AgentExecution, error)
	GetGoalAchievedImpactedIdsLast30Days(ctx context.Context, agentID string) ([]string, error)
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
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	if executionRecord.AgentID == nil {
		span.LogFields(log.Object("executionRecord", executionRecord))
		err := errors.New("agent or executionID missing")
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
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
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

	fields := map[string]any{
		"status":        enum.AgentExecutionCompleted.String(),
		"completed_at":  utils.NowPtr(),
		"next_retry_at": nil,
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
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.String("executionID", executionID))

	err := f.gormDb.
		Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Updates(map[string]any{
			"status":        enum.AgentExecutionCompleted.String(),
			"error_message": nil,
			"completed_at":  utils.Now(),
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
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.String("executionID", executionID))

	err := f.gormDb.
		Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Where("status <> ?", enum.AgentExecutionCompleted.String()).
		Updates(map[string]any{
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
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
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
		Updates(map[string]any{
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
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

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
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.String("executionID", executionID))
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

func (f *agentExecutionRepository) ScheduleRetry(ctx context.Context, executionID string, inputError error, stateData map[string]any) (*time.Time, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.ScheduleRetry")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.String("executionID", executionID))
	if inputError == nil {
		span.LogFields(log.String("inputError", "nil"))
	} else {
		span.LogFields(log.String("inputError", inputError.Error()))
	}

	execution := &postgres_entity.AgentExecution{}
	result := f.gormDb.First(execution, "id = ?", executionID)
	if result.Error != nil {
		return nil, result.Error
	}

	// Prepare update fields
	updateFields := map[string]any{
		"state_data": stateData,
	}

	// Increment retry count and set next retry time if it's not the first scheduling
	if execution.NextRetryAt != nil {
		updateFields["retry_count"] = execution.RetryCount + 1
	}

	// Check retry conditions and set appropriate fields
	if execution.RetryCount != 0 && execution.RetryCount >= execution.MaxRetries {
		updateFields["status"] = enum.AgentExecutionError.String()
		if inputError != nil {
			updateFields["error_message"] = inputError.Error()
		}
		updateFields["next_retry_at"] = nil
	} else {
		backoff := utils.CalculateExponentialBackoffDelay(execution.RetryCount, utils.BackoffConfig{
			InitialDelay: 30 * time.Second,
			MaxDelay:     12 * time.Hour,
			Factor:       2.0,
			Jitter:       0.12,
		})
		nextRetry := utils.Now().Add(backoff)
		updateFields["next_retry_at"] = nextRetry
		if inputError != nil {
			updateFields["error_message"] = inputError.Error()
		}
		updateFields["status"] = enum.AgentExecutionRetrying.String()
	}

	// Update the record
	result = f.gormDb.Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Updates(updateFields)
	if result.Error != nil {
		return nil, result.Error
	}

	// reload the record
	result = f.gormDb.First(execution, "id = ?", executionID)
	if result.Error != nil {
		return nil, result.Error
	}

	return execution.NextRetryAt, nil
}

func (f *agentExecutionRepository) SaveAsyncState(ctx context.Context, executionID string, currentStep string, stateData map[string]any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.SaveAsyncState")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.String("executionID", executionID))

	return f.gormDb.Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Updates(map[string]any{
			"current_step": currentStep,
			"state_data":   stateData,
			"status":       enum.AgentExecutionPending,
		}).Error
}

func (f *agentExecutionRepository) CompleteStep(ctx context.Context, executionID string, step string, result map[string]any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.CompleteStep")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.String("executionID", executionID))

	// First get the current execution to check/initialize checkpoints
	execution := &postgres_entity.AgentExecution{}
	err := f.gormDb.First(execution, "id = ?", executionID).Error
	if err != nil {
		return err
	}

	// Initialize or update checkpoints
	checkpoints := execution.Checkpoints
	if checkpoints == nil {
		checkpoints = make(map[string]any)
	}
	checkpoints[step] = result

	// Update only the specific fields we want to change
	return f.gormDb.Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Updates(map[string]any{
			"checkpoints":  checkpoints,
			"current_step": "",
			"state_data":   nil,
		}).Error
}

func (f *agentExecutionRepository) GoalAchieved(ctx context.Context, executionID string, goalAchieved bool, impactedId *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.GoalAchieved")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	tracing.TagEntity(span, executionID)
	span.LogFields(log.Bool("goalAchieved", goalAchieved))

	if executionID == "" {
		err := errors.New("ID is missing")
		tracing.TraceErr(span, err)
		return err
	}

	fieldsToUpdate := map[string]any{
		"goal_achieved": goalAchieved,
		"completed_at":  utils.NowPtr(),
	}
	if impactedId != nil {
		fieldsToUpdate["impacted_id"] = *impactedId
	}

	err := f.gormDb.
		Model(&postgres_entity.AgentExecution{}).
		Where("id = ?", executionID).
		Updates(fieldsToUpdate).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *agentExecutionRepository) GetGoalAchievedCountLast30Days(ctx context.Context, agentID string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.GetGoalAchievedCountLast30Days")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	var count int64
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	err := r.gormDb.Model(&postgres_entity.AgentExecution{}).
		Where("agent_id = ? AND goal_achieved = true AND updated_at >= ?", agentID, thirtyDaysAgo).
		Count(&count).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return 0, err
	}

	return count, nil
}

func (f *agentExecutionRepository) GetExecutionsForRetry(ctx context.Context, limit int) ([]postgres_entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.GetExecutionsForRetry")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.Int("limit", limit))

	var executions []postgres_entity.AgentExecution
	err := f.gormDb.
		Where("status = ?", enum.AgentExecutionRetrying.String()).
		Where("next_retry_at IS NULL OR next_retry_at <= ?", utils.Now()).
		Where("retry_count <= max_retries").
		Limit(limit).
		Find(&executions).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return executions, nil
}

func (f *agentExecutionRepository) GetGoalAchievedImpactedIdsLast30Days(ctx context.Context, agentID string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentExecutionRepository.GetGoalAchievedImpactedIdsLast30Days")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	var impactedIds []string
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	err := f.gormDb.Model(&postgres_entity.AgentExecution{}).
		Where("agent_id = ? AND goal_achieved = true AND updated_at >= ? AND impacted_id IS NOT NULL", agentID, thirtyDaysAgo).
		Pluck("impacted_id", &impactedIds).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	impactedIds = utils.RemoveEmpties(impactedIds)
	impactedIds = utils.RemoveDuplicates(impactedIds)

	span.LogFields(log.Int("result.count", len(impactedIds)))
	return impactedIds, nil
}
