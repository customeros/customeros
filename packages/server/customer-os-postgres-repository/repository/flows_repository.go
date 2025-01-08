package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowsRepository interface {
	Create(ctx context.Context, flowRecord entity.Flows) (*entity.Flows, error)
	Find(ctx context.Context, flowRecord entity.Flows) (*entity.Flows, error)
	FindAll(ctx context.Context, flowRecord entity.Flows) (*[]entity.Flows, error)
	Update(ctx context.Context, flowRecord entity.Flows) (*entity.Flows, error)
}

type flowRepository struct {
	gormDb *gorm.DB
}

func NewFlowsRepository(gormDb *gorm.DB) FlowsRepository {
	return &flowRepository{gormDb: gormDb}
}

func (f *flowRepository) Create(ctx context.Context, flowRecord entity.Flows) (*entity.Flows, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowsRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	flowRecord.ID = utils.GenerateNanoIdWithPrefix("flow", 16)

	var created entity.Flows
	err := f.gormDb.Create(&flowRecord).Scan(&created).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &flowRecord, nil
}

func (f *flowRepository) FindAll(ctx context.Context, flowRecord entity.Flows) (*[]entity.Flows, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowsRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var flows []entity.Flows
	query := f.gormDb.Where("is_active = true")

	// Add additional filters based on non-zero fields in flowRecord
	if flowRecord != (entity.Flows{}) {
		query = query.Where(&flowRecord)
	}

	err := query.Find(&flows).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &flows, nil
}

func (f *flowRepository) Find(ctx context.Context, flowRecord entity.Flows) (*entity.Flows, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowsRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var flow entity.Flows
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

func (f *flowRepository) Update(ctx context.Context, flowRecord entity.Flows) (*entity.Flows, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowsRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowRecord.ID == "" {
		err := errors.New("flowID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedFlows entity.Flows
	err := f.gormDb.
		Model(&flowRecord).
		Updates(&flowRecord).
		First(&updatedFlows).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedFlows, nil
}
