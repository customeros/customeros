package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type FlowRepository interface {
	GetList(ctx context.Context) ([]*entity.Flow, error)
	GetById(ctx context.Context, id string) (*entity.Flow, error)

	Store(ctx context.Context, entity *entity.Flow) (*entity.Flow, error)
}

type flowRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewFlowRepository(gormDb *gorm.DB) FlowRepository {
	return &flowRepositoryImpl{gormDb: gormDb}
}

func (r flowRepositoryImpl) GetList(ctx context.Context) ([]*entity.Flow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowRepository.Get")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	var result []*entity.Flow
	err := r.gormDb.
		Where("tenant = ?", common.GetTenantFromContext(ctx)).
		Find(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(result)))

	return result, nil
}

func (r flowRepositoryImpl) GetById(ctx context.Context, id string) (*entity.Flow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.String("id", id))

	var result entity.Flow
	err := r.gormDb.
		Where("tenant = ? and id = ?", common.GetTenantFromContext(ctx), id).
		First(&result).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			span.LogFields(tracingLog.Bool("result.found", false))
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Bool("result.found", true))

	return &result, nil
}

func (repo *flowRepositoryImpl) Store(ctx context.Context, entity *entity.Flow) (*entity.Flow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowRepository.Store")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.Object("entity", entity))

	err := repo.gormDb.Save(entity).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.String("entity.id", entity.ID))

	return entity, nil
}
