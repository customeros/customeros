package postgres_repository

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type FlowNodeRepository interface {
	Create(ctx context.Context, flowNode postgres_entity.FlowNode) (*postgres_entity.FlowNode, error)
	FindAll(ctx context.Context, flowNode postgres_entity.FlowNode) ([]postgres_entity.FlowNode, error)
	Find(ctx context.Context, flowNode postgres_entity.FlowNode) (*postgres_entity.FlowNode, error)
	Update(ctx context.Context, flowNode postgres_entity.FlowNode) (*postgres_entity.FlowNode, error)
}

type flowNodeRepository struct {
	gormDb *gorm.DB
}

func NewFlowNodeRepository(gormDb *gorm.DB) FlowNodeRepository {
	return &flowNodeRepository{gormDb: gormDb}
}

func (f *flowNodeRepository) Create(ctx context.Context, flowNode postgres_entity.FlowNode) (*postgres_entity.FlowNode, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowNodeRepository.Create")
	defer spans.Finish()

	flowNode.ID = utils.GenerateNanoIdWithPrefix("node", 16)

	var created postgres_entity.FlowNode
	err := f.gormDb.Create(&flowNode).Scan(&created).Error
	if err != nil {
		return nil, err
	}
	return &flowNode, nil
}

func (f *flowNodeRepository) FindAll(ctx context.Context, flowNode postgres_entity.FlowNode) ([]postgres_entity.FlowNode, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowRepository.FindAll")
	defer spans.Finish()

	var nodes []postgres_entity.FlowNode
	err := f.gormDb.
		Where(&flowNode).
		Find(&nodes).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return nodes, nil
}

func (f *flowNodeRepository) Find(ctx context.Context, flowNode postgres_entity.FlowNode) (*postgres_entity.FlowNode, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowNodeRepository.Find")
	defer spans.Finish()

	var node postgres_entity.FlowNode
	err := f.gormDb.
		Where(&flowNode).
		First(&node).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}
	return &node, nil
}

func (f *flowNodeRepository) Update(ctx context.Context, flowNode postgres_entity.FlowNode) (*postgres_entity.FlowNode, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "FlowNodeRepository.Update")
	defer spans.Finish()

	if flowNode.ID == "" {
		err := errors.New("flow node ID is missing")
		spans.TraceError(err)
		return nil, err
	}

	var updatedNode postgres_entity.FlowNode
	err := f.gormDb.
		Model(&postgres_entity.FlowNode{}).
		Where("id = ?", flowNode.ID).
		Updates(flowNode).
		First(&updatedNode, "id = ?", flowNode.ID).
		Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogObjectAsJson("payload", updatedNode)
	return &updatedNode, nil
}
