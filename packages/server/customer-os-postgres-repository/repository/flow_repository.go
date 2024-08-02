package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type FlowRepository interface {
	Count(ctx context.Context, tenant string) (int64, error)
	Get(ctx context.Context, tenant string, page, limit int) ([]*entity.Flow, error)
	GetById(ctx context.Context, tenant, id string) (*entity.Flow, error)

	Store(ctx context.Context, entity *entity.Flow) (*entity.Flow, error)
	Delete(ctx context.Context, tenant, id string) error
}

type flowRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewFlowRepository(gormDb *gorm.DB) FlowRepository {
	return &flowRepositoryImpl{gormDb: gormDb}
}

func (r flowRepositoryImpl) Count(ctx context.Context, tenant string) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowRepository.Count")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	var result int64
	err := r.gormDb.
		Model(entity.Flow{}).
		Count(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return 0, err
	}

	span.LogFields(tracingLog.Int64("result.count", result))

	return result, nil
}

func (r flowRepositoryImpl) Get(ctx context.Context, tenant string, page, limit int) ([]*entity.Flow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowRepository.Get")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	span.LogFields(tracingLog.Int("page", page), tracingLog.Int("limit", limit))

	var result []*entity.Flow
	err := r.gormDb.
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(result)))

	return result, nil
}

func (r flowRepositoryImpl) GetById(ctx context.Context, tenant, id string) (*entity.Flow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowRepository.GetById")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	span.LogFields(tracingLog.String("id", id))

	var result entity.Flow
	err := r.gormDb.
		Where("id = ?", id).
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

	span.SetTag(tracing.SpanTagTenant, entity.Tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.Object("entity", entity))

	err := repo.gormDb.Save(entity).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.String("entity.id", entity.ID))

	return entity, nil
}

func (repo *flowRepositoryImpl) Delete(ctx context.Context, tenant, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowRepository.Delete")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("id", id))

	err := repo.gormDb.Delete(&entity.Flow{}, "id = ?", id).Error

	if err != nil {
		span.LogFields(tracingLog.Bool("entity.deleted", false))
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(tracingLog.Bool("entity.deleted", true))

	return nil
}
