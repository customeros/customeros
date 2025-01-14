package repository

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowNodeRepository interface {
	Create(ctx context.Context, flowNode entity.FlowNode) (*entity.FlowNode, error)
	FindAll(ctx context.Context, flowNode entity.FlowNode) ([]entity.FlowNode, error)
	Find(ctx context.Context, flowNode entity.FlowNode) (*entity.FlowNode, error)
	Update(ctx context.Context, flowNode entity.FlowNode) (*entity.FlowNode, error)
}

type flowNodeRepository struct {
	gormDb *gorm.DB
}

func NewFlowNodeRepository(gormDb *gorm.DB) FlowNodeRepository {
	return &flowNodeRepository{gormDb: gormDb}
}

func (f *flowNodeRepository) Create(ctx context.Context, flowNode entity.FlowNode) (*entity.FlowNode, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowNodeRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	flowNode.ID = utils.GenerateNanoIdWithPrefix("node", 16)

	var created entity.FlowNode
	err := f.gormDb.Create(&flowNode).Scan(&created).Error
	if err != nil {
		return nil, err
	}
	return &flowNode, nil
}

func (f *flowNodeRepository) FindAll(ctx context.Context, flowNode entity.FlowNode) ([]entity.FlowNode, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowRepository.FindAll")
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
	return nodes, nil
}

func (f *flowNodeRepository) Find(ctx context.Context, flowNode entity.FlowNode) (*entity.FlowNode, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowNodeRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var node entity.FlowNode
	err := f.gormDb.
		Where(&flowNode).
		First(&node).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &node, nil
}

func (f *flowNodeRepository) Update(ctx context.Context, flowNode entity.FlowNode) (*entity.FlowNode, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowNodeRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowNode.ID == "" {
		err := errors.New("flow node ID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedNode entity.FlowNode
	err := f.gormDb.
		Model(&entity.FlowNode{}).
		Where("id = ?", flowNode.ID).
		Updates(flowNode).
		First(&updatedNode, "id = ?", flowNode.ID).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	tracing.LogObjectAsJson(span, "payload", updatedNode)
	return &updatedNode, nil
}
