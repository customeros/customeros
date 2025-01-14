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

type FlowEdgeRepository interface {
	Create(ctx context.Context, flowEdge entity.FlowEdge) (*entity.FlowEdge, error)
	FindAll(ctx context.Context, flowEdge entity.FlowEdge) ([]entity.FlowEdge, error)
	Find(ctx context.Context, flowEdge entity.FlowEdge) (*entity.FlowEdge, error)
	Update(ctx context.Context, flowEdge entity.FlowEdge) (*entity.FlowEdge, error)
}

type flowEdgeRepository struct {
	gormDb *gorm.DB
}

func NewFlowEdgeRepository(gormDb *gorm.DB) FlowEdgeRepository {
	return &flowEdgeRepository{gormDb: gormDb}
}

func (f *flowEdgeRepository) Create(ctx context.Context, flowEdge entity.FlowEdge) (*entity.FlowEdge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowEdgeRepository.CreateFlowEdge")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	flowEdge.ID = utils.GenerateNanoIdWithPrefix("edge", 16)

	var created entity.FlowEdge
	err := f.gormDb.Create(&flowEdge).Scan(&created).Error
	if err != nil {
		return nil, err
	}
	return &flowEdge, nil
}

func (f *flowEdgeRepository) FindAll(ctx context.Context, flowEdge entity.FlowEdge) ([]entity.FlowEdge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowEdgeRepository.FindAll")
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
	return results, nil
}

func (f *flowEdgeRepository) Find(ctx context.Context, flowEdge entity.FlowEdge) (*entity.FlowEdge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowEdgeRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var result entity.FlowEdge
	err := f.gormDb.
		Where(&flowEdge).
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &result, nil
}

func (f *flowEdgeRepository) Update(ctx context.Context, flowEdge entity.FlowEdge) (*entity.FlowEdge, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowEdgeRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowEdge.ID == "" {
		err := errors.New("flow edge ID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedEdge entity.FlowEdge
	err := f.gormDb.
		Model(&entity.FlowEdge{}).
		Where("id = ?", flowEdge.ID).
		Updates(flowEdge).
		First(&updatedEdge, "id = ?", flowEdge.ID).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedEdge, nil
}
