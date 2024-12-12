package repository

import (
	"context"
	"fmt"

	nanoid "github.com/matoous/go-nanoid/v2"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowRepository interface {
	CreateFlow(ctx context.Context, flow entity.Flow) (entity.Flow, error)
	GetAllFlowsForTenant(ctx context.Context) ([]entity.Flow, error)
}

type flowRepository struct {
	gormDb *gorm.DB
}

func NewFlowRepository(gormDb *gorm.DB) FlowRepository {
	return &flowRepository{gormDb: gormDb}
}

func (f *flowRepository) CreateFlow(ctx context.Context, flow entity.Flow) (entity.Flow, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "FlowRepository.CreateFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	flow.ID = fmt.Sprintf("flow-%s", nanoid.Must())

	err := f.gormDb.Create(&flow).Error
	if err != nil {
		return entity.Flow{}, err
	}
	return flow, nil
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
