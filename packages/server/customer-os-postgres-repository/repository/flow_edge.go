package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type FlowEdgeRepository interface {
	Create(ctx context.Context, flowEdge postgres_entity.FlowEdge) (*postgres_entity.FlowEdge, error)
	FindAll(ctx context.Context, flowEdge postgres_entity.FlowEdge) ([]postgres_entity.FlowEdge, error)
	Find(ctx context.Context, flowEdge postgres_entity.FlowEdge) (*postgres_entity.FlowEdge, error)
	Update(ctx context.Context, flowEdge postgres_entity.FlowEdge) (*postgres_entity.FlowEdge, error)
}

type flowEdgeRepository struct {
	gormDb *gorm.DB
}

func NewFlowEdgeRepository(gormDb *gorm.DB) FlowEdgeRepository {
	return &flowEdgeRepository{gormDb: gormDb}
}

func (f *flowEdgeRepository) Create(ctx context.Context, flowEdge postgres_entity.FlowEdge) (*postgres_entity.FlowEdge, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowEdgeRepository.CreateFlowEdge")
	defer spans.Finish()

	flowEdge.ID = utils.GenerateNanoIdWithPrefix("edge", 16)

	var created postgres_entity.FlowEdge
	err := f.gormDb.Create(&flowEdge).Scan(&created).Error
	if err != nil {
		return nil, err
	}
	return &flowEdge, nil
}

func (f *flowEdgeRepository) FindAll(ctx context.Context, flowEdge postgres_entity.FlowEdge) ([]postgres_entity.FlowEdge, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowEdgeRepository.FindAll")
	defer spans.Finish()

	var results []postgres_entity.FlowEdge
	err := f.gormDb.
		Where(&flowEdge).
		Find(&results).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return results, nil
}

func (f *flowEdgeRepository) Find(ctx context.Context, flowEdge postgres_entity.FlowEdge) (*postgres_entity.FlowEdge, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowEdgeRepository.Find")
	defer spans.Finish()

	var result postgres_entity.FlowEdge
	err := f.gormDb.
		Where(&flowEdge).
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	return &result, nil
}

func (f *flowEdgeRepository) Update(ctx context.Context, flowEdge postgres_entity.FlowEdge) (*postgres_entity.FlowEdge, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowEdgeRepository.Update")
	defer spans.Finish()

	if flowEdge.ID == "" {
		err := errors.New("flow edge ID is missing")
		spans.TraceError(err)
		return nil, err
	}

	var updatedEdge postgres_entity.FlowEdge
	err := f.gormDb.
		Model(&postgres_entity.FlowEdge{}).
		Where("id = ?", flowEdge.ID).
		Updates(flowEdge).
		First(&updatedEdge, "id = ?", flowEdge.ID).
		Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return &updatedEdge, nil
}
