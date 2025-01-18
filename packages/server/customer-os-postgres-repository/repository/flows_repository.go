package postgres_repository

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowsRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	flowRecord.ID = utils.GenerateNanoIdWithPrefix("flow", 16)

	var created postgres_entity.Flows
	err := f.gormDb.Create(&flowRecord).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &flowRecord, nil
}

func (f *flowRepository) FindAll(ctx context.Context, flowRecord postgres_entity.Flows) ([]postgres_entity.Flows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowsRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var flows []postgres_entity.Flows
	query := f.gormDb.Where("is_active = true")

	// Add additional filters based on non-zero fields in flowRecord
	if flowRecord != (postgres_entity.Flows{}) {
		query = query.Where(&flowRecord)
	}

	err := query.Find(&flows).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return flows, nil
}

func (f *flowRepository) Find(ctx context.Context, flowRecord postgres_entity.Flows) (*postgres_entity.Flows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowsRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var flow postgres_entity.Flows
	err := f.gormDb.
		Where(&flowRecord).
		Where("is_active = true").
		First(&flow).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &flow, nil
}

func (f *flowRepository) Update(ctx context.Context, flowRecord postgres_entity.Flows) (*postgres_entity.Flows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowsRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowRecord.ID == "" {
		err := errors.New("flow ID is missing")
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedFlows, nil
}
