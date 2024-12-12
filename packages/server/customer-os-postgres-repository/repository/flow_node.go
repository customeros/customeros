package repository

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowNodeRepository interface {
	CreateFlowNode(ctx context.Context, flowNode entity.FlowNode) (entity.FlowNode, error)
	GetAllNodesForFlow(ctx context.Context, flowID string) ([]entity.FlowNode, error)
	GetNodeById(ctx context.Context, nodeId string) (entity.FlowNode, error)
}

type flowNodeRepository struct {
	gormDb *gorm.DB
}

const (
	DefaultNodeX float64 = 100
	DefaultNodeY float64 = 100
)

func NewFlowNodeRepository(gormDb *gorm.DB) FlowNodeRepository {
	return &flowNodeRepository{gormDb: gormDb}
}

func (f *flowNodeRepository) CreateFlowNode(ctx context.Context, flowNode entity.FlowNode) (entity.FlowNode, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowNodeRepository.CreateFlowNode")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	flowNode.ID = utils.GenerateNanoIdWithPrefix("node")

	x := DefaultNodeX
	y := DefaultNodeY

	if int(*flowNode.PositionX) == 0 {
		flowNode.PositionX = &x
	}
	if int(*flowNode.PositionY) == 0 {
		flowNode.PositionY = &y
	}

	err := f.gormDb.Create(&flowNode).Error
	if err != nil {
		return entity.FlowNode{}, err
	}
	return flowNode, nil
}

func (f *flowNodeRepository) GetAllNodesForFlow(ctx context.Context, flowID string) ([]entity.FlowNode, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.GetAllNodesForFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var nodes []entity.FlowNode
	err := f.gormDb.
		Where("tenant = ? AND flow_id = ?", common.GetTenantFromContext(ctx), flowID).
		Find(&nodes).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nodes, err
	}
	return nodes, nil
}

func (f *flowNodeRepository) GetNodeById(ctx context.Context, nodeId string) (entity.FlowNode, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.GetNodeById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var node entity.FlowNode
	err := f.gormDb.
		Where("tenant = ? AND id = ?", common.GetTenantFromContext(ctx), nodeId).
		First(&node).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return node, err
	}
	return node, nil
}
