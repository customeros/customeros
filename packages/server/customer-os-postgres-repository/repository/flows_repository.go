package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type FlowsRepository interface {
	Create(ctx context.Context, flowRecord postgres_entity.Flows) (*postgres_entity.Flows, error)
	Find(ctx context.Context, flowRecord postgres_entity.Flows) (*postgres_entity.Flows, error)
	FindAll(ctx context.Context, flowRecord postgres_entity.Flows) ([]postgres_entity.Flows, error)
	Update(ctx context.Context, flowRecord postgres_entity.Flows) (*postgres_entity.Flows, error)
}

type flowRepository struct {
	gormDb *gorm.DB
}

func NewFlowsRepository(gormDb *gorm.DB) FlowsRepository {
	return &flowRepository{gormDb: gormDb}
}

func (f *flowRepository) Create(ctx context.Context, flowRecord postgres_entity.Flows) (*postgres_entity.Flows, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowsRepository.Create")
	defer spans.Finish()

	flowRecord.ID = utils.GenerateNanoIdWithPrefix("flow", 16)

	var created postgres_entity.Flows
	err := f.gormDb.Create(&flowRecord).Scan(&created).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return &flowRecord, nil
}

func (f *flowRepository) FindAll(ctx context.Context, flowRecord postgres_entity.Flows) ([]postgres_entity.Flows, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowsRepository.FindAll")
	defer spans.Finish()

	var flows []postgres_entity.Flows
	query := f.gormDb.Where("is_active = true")

	// Add additional filters based on non-zero fields in flowRecord
	if flowRecord != (postgres_entity.Flows{}) {
		query = query.Where(&flowRecord)
	}

	err := query.Find(&flows).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return flows, nil
}

func (f *flowRepository) Find(ctx context.Context, flowRecord postgres_entity.Flows) (*postgres_entity.Flows, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowsRepository.Find")
	defer spans.Finish()

	var flow postgres_entity.Flows
	err := f.gormDb.
		Where(&flowRecord).
		Where("is_active = true").
		First(&flow).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}
	return &flow, nil
}

func (f *flowRepository) Update(ctx context.Context, flowRecord postgres_entity.Flows) (*postgres_entity.Flows, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowsRepository.Update")
	defer spans.Finish()

	if flowRecord.ID == "" {
		err := errors.New("flow ID is missing")
		spans.TraceError(err)
		return nil, err
	}

	var updatedFlows postgres_entity.Flows
	err := f.gormDb.
		Model(&postgres_entity.Flows{}).
		Where("id = ?", flowRecord.ID).
		Updates(flowRecord).
		First(&updatedFlows, "id = ?", flowRecord.ID).
		Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &updatedFlows, nil
}
