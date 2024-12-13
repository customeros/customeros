package repository

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowEdgeRepository interface {
	Create(ctx context.Context, flowEdge entity.FlowEdge) (*entity.FlowEdge, error)
	FindAll(ctx context.Context, flowEdge entity.FlowEdge) (*[]entity.FlowEdge, error)
}

type flowEdgeRepository struct {
	gormDb *gorm.DB
}

func NewFlowEdgeRepository(gormDb *gorm.DB) FlowEdgeRepository {
	return &flowEdgeRepository{gormDb: gormDb}
}

func (f *flowEdgeRepository) Create(ctx context.Context, flowEdge entity.FlowEdge) (*entity.FlowEdge, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowEdgeRepository.CreateFlowEdge")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	flowEdge.ID = utils.GenerateNanoIdWithPrefix("edge")

	err := f.gormDb.Create(&flowEdge).Error
	if err != nil {
		return nil, err
	}
	return &flowEdge, nil
}

func (f *flowEdgeRepository) FindAll(ctx context.Context, flowEdge entity.FlowEdge) (*[]entity.FlowEdge, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowEdgeRepository.FindEdges")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var results []entity.FlowEdge
	err := f.gormDb.
		Where(&flowEdge).
		Find(&results).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &results, nil
}
