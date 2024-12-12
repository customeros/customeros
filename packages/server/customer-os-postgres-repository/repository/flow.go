package repository

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowRepository interface {
	CreateFlow(ctx context.Context, flowRecord entity.Flow) (entity.Flow, error)
	GetAllFlowsForTenant(ctx context.Context) ([]entity.Flow, error)
	GetFlowByID(ctx context.Context, flowID string) (*entity.Flow, error)
}

type flowRepository struct {
	gormDb *gorm.DB
}

func NewFlowRepository(gormDb *gorm.DB) FlowRepository {
	return &flowRepository{gormDb: gormDb}
}

func (f *flowRepository) CreateFlow(ctx context.Context, flowRecord entity.Flow) (entity.Flow, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.CreateFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	flowRecord.ID = utils.GenerateNanoIdWithPrefix("flow")

	err := f.gormDb.Create(&flowRecord).Error
	if err != nil {
		return entity.Flow{}, err
	}
	return flowRecord, nil
}

func (f *flowRepository) GetAllFlowsForTenant(ctx context.Context) ([]entity.Flow, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.GetAllFlows")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var flows []entity.Flow
	err := f.gormDb.
		Where("tenant = ?", common.GetTenantFromContext(ctx)).
		Find(&flows).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return flows, err
	}
	return flows, nil
}

func (f *flowRepository) GetFlowByID(ctx context.Context, flowID string) (*entity.Flow, error) {

	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.GetFlowByID")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var flow entity.Flow
	err := f.gormDb.
		Where("tenant = ? AND id = ?", common.GetTenantFromContext(ctx), flowID).
		First(&flow).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return &flow, err
	}
	return &flow, nil

}
