package repository

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowEdgeRepository interface {
	CreateFlowEdge(ctx context.Context, flowEdge entity.FlowEdge) (entity.FlowEdge, error)
	GetAllEdgesForFlow(ctx context.Context, FlowID string) ([]entity.FlowEdge, error)
}

type flowEdgeRepository struct {
	gormDb *gorm.DB
}

func NewFlowEdgeRepository(gormDb *gorm.DB) FlowEdgeRepository {
	return &flowEdgeRepository{gormDb: gormDb}
}

func (f *flowEdgeRepository) CreateFlowEdge(ctx context.Context, flowEdge entity.FlowEdge) (entity.FlowEdge, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowEdgeRepository.CreateFlowEdge")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	flowEdge.ID = utils.GenerateNanoIdWithPrefix("node")

	err := f.gormDb.Create(&flowEdge).Error
	if err != nil {
		return entity.FlowEdge{}, err
	}
	return flowEdge, nil
}

func (f *flowEdgeRepository) GetAllEdgesForFlow(ctx context.Context, flowID string) ([]entity.FlowEdge, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.GetAllEdgesForFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var nodes []entity.FlowEdge
	err := f.gormDb.
		Where("tenant = ? AND flow_id = ?", common.GetTenantFromContext(ctx), flowID).
		Find(&nodes).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nodes, err
	}
	return nodes, nil
}
