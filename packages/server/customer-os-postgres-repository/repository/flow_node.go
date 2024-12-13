package repository

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowNodeRepository interface {
	Create(ctx context.Context, flowNode entity.FlowNode) (*entity.FlowNode, error)
	FindAll(ctx context.Context, flowNode entity.FlowNode) (*[]entity.FlowNode, error)
	Find(ctx context.Context, flowNode entity.FlowNode) (*entity.FlowNode, error)
}

type flowNodeRepository struct {
	gormDb *gorm.DB
}

func NewFlowNodeRepository(gormDb *gorm.DB) FlowNodeRepository {
	return &flowNodeRepository{gormDb: gormDb}
}

func (f *flowNodeRepository) Create(ctx context.Context, flowNode entity.FlowNode) (*entity.FlowNode, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowNodeRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	flowNode.ID = utils.GenerateNanoIdWithPrefix("node")

	err := f.gormDb.Create(&flowNode).Error
	if err != nil {
		return nil, err
	}
	return &flowNode, nil
}

func (f *flowNodeRepository) FindAll(ctx context.Context, flowNode entity.FlowNode) (*[]entity.FlowNode, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var nodes []entity.FlowNode
	err := f.gormDb.
		Where(&flowNode).
		Find(&nodes).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &nodes, nil
}

func (f *flowNodeRepository) Find(ctx context.Context, flowNode entity.FlowNode) (*entity.FlowNode, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var node entity.FlowNode
	err := f.gormDb.
		Where(&flowNode).
		First(&node).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &node, nil
}
