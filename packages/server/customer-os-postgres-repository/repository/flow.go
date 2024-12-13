package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowRepository interface {
	Create(ctx context.Context, flowRecord entity.Flow) (*entity.Flow, error)
	Find(ctx context.Context, flowRecord entity.Flow) (*entity.Flow, error)
	FindAll(ctx context.Context, flowRecord entity.Flow) (*[]entity.Flow, error)
	Update(ctx context.Context, flowRecord entity.Flow) (*entity.Flow, error)
}

type flowRepository struct {
	gormDb *gorm.DB
}

func NewFlowRepository(gormDb *gorm.DB) FlowRepository {
	return &flowRepository{gormDb: gormDb}
}

func (f *flowRepository) Create(ctx context.Context, flowRecord entity.Flow) (*entity.Flow, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowRecord.Tenant == "" {
		flowRecord.Tenant = common.GetTenantFromContext(ctx)
		if flowRecord.Tenant == "" {
			err := errors.New("Tenant not set in context")
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	flowRecord.ID = utils.GenerateNanoIdWithPrefix("flow", 16)

	err := f.gormDb.Create(&flowRecord).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &flowRecord, nil
}

func (f *flowRepository) FindAll(ctx context.Context, flowRecord entity.Flow) (*[]entity.Flow, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.FindAll")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowRecord.Tenant == "" {
		flowRecord.Tenant = common.GetTenantFromContext(ctx)
		if flowRecord.Tenant == "" {
			err := errors.New("Tenant not set in context")
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	var flows []entity.Flow
	query := f.gormDb.Where("tenant = ? AND status != ?", flowRecord.Tenant, enum.FlowStatusArchived)

	// Add additional filters based on non-zero fields in flowRecord
	if flowRecord != (entity.Flow{}) {
		query = query.Where(&flowRecord)
	}

	err := query.Find(&flows).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &flows, nil
}

func (f *flowRepository) Find(ctx context.Context, flowRecord entity.Flow) (*entity.Flow, error) {

	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowRecord.Tenant == "" {
		flowRecord.Tenant = common.GetTenantFromContext(ctx)
		if flowRecord.Tenant == "" {
			err := errors.New("Tenant not set in context")
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	var flow entity.Flow
	err := f.gormDb.
		Where(&flowRecord).
		Where("status != ?", enum.FlowStatusArchived).
		First(&flow).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &flow, nil

}

func (f *flowRepository) Update(ctx context.Context, flowRecord entity.Flow) (*entity.Flow, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowRecord.Tenant == "" {
		flowRecord.Tenant = common.GetTenantFromContext(ctx)
		if flowRecord.Tenant == "" {
			err := errors.New("tenant not set in context")
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	if flowRecord.ID == "" {
		err := errors.New("flowID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var updatedFlow entity.Flow
	err := f.gormDb.
		Model(&flowRecord).
		Updates(&flowRecord).
		First(&updatedFlow).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedFlow, nil
}
