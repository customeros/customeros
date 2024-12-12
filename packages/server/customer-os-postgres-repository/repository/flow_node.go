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
	GetAllNodesForFlow(ctx context.Context, FlowID string) ([]entity.FlowNode, error)
}

type flowNodeRepository struct {
	gormDb *gorm.DB
}

func NewFlowNodeRepository(gormDb *gorm.DB) FlowNodeRepository {
	return &flowNodeRepository{gormDb: gormDb}
}

func (f *flowNodeRepository) CreateFlowNode(ctx context.Context, flowNode entity.FlowNode) (entity.FlowNode, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowNodeRepository.CreateFlowNode")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	flowNode.ID = utils.GenerateNanoIdWithPrefix("node")

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
