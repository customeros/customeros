package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type FlowExecutionRepository interface {
	Create(ctx context.Context, executionRecord postgres_entity.FlowExecution) (*postgres_entity.FlowExecution, error)
	Find(ctx context.Context, executionRecord postgres_entity.FlowExecution) (*postgres_entity.FlowExecution, error)
	Update(ctx context.Context, executionRecord postgres_entity.FlowExecution) (*postgres_entity.FlowExecution, error)
}

type flowExecutionRepository struct {
	gormDb *gorm.DB
}

func NewFlowExecutionRepository(gormDb *gorm.DB) FlowExecutionRepository {
	return &flowExecutionRepository{gormDb: gormDb}
}

func (f *flowExecutionRepository) Create(ctx context.Context, executionRecord postgres_entity.FlowExecution) (*postgres_entity.FlowExecution, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowExecutionRepository.Create")
	defer spans.Finish()

	err := f.gormDb.Create(&executionRecord).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &executionRecord, nil
}

func (f *flowExecutionRepository) Find(ctx context.Context, executionRecord postgres_entity.FlowExecution) (*postgres_entity.FlowExecution, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowExecutionRepository.Find")
	defer spans.Finish()

	var foundRecord postgres_entity.FlowExecution
	err := f.gormDb.
		Where(&executionRecord).
		First(&foundRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}
	return &foundRecord, nil
}

func (f *flowExecutionRepository) FindAll(ctx context.Context, executionRecord postgres_entity.FlowExecution) (*[]postgres_entity.FlowExecution, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowExecutionRepository.FindAll")
	defer spans.Finish()

	var records []postgres_entity.FlowExecution
	err := f.gormDb.
		Where(&executionRecord).
		Find(&records).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return &records, nil
}

func (f *flowExecutionRepository) Update(ctx context.Context, executionRecord postgres_entity.FlowExecution) (*postgres_entity.FlowExecution, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowExecutionRepository.Update")
	defer spans.Finish()

	if executionRecord.ID == "" {
		err := errors.New("flow execution ID is missing")
		spans.TraceError(err)
		return nil, err
	}

	var updatedRecord postgres_entity.FlowExecution
	err := f.gormDb.
		Model(&postgres_entity.FlowExecution{}).
		Where("id = ?", executionRecord.ID).
		Updates(executionRecord).
		First(&updatedRecord, "id = ?", executionRecord.ID).
		Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &updatedRecord, nil
}
