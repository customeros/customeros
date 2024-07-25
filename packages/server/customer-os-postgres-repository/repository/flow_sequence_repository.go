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

type FlowSequenceRepository interface {
	Get(ctx context.Context, tenant, flowId string) ([]*entity.FlowSequence, error)
	GetById(ctx context.Context, tenant, id string) (*entity.FlowSequence, error)

	Store(ctx context.Context, tenant string, entity *entity.FlowSequence) (*entity.FlowSequence, error)
	Delete(ctx context.Context, tenant, id string) error
}

type flowSequenceRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewFlowSequenceRepository(gormDb *gorm.DB) FlowSequenceRepository {
	return &flowSequenceRepositoryImpl{gormDb: gormDb}
}

func (r flowSequenceRepositoryImpl) Get(ctx context.Context, tenant, flowId string) ([]*entity.FlowSequence, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceRepository.Get")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	var result []*entity.FlowSequence
	err := r.gormDb.
		Where("flow_id = ?", flowId).
		Find(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(result)))

	return result, nil
}

func (r flowSequenceRepositoryImpl) GetById(ctx context.Context, tenant, id string) (*entity.FlowSequence, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceRepository.GetById")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	span.LogFields(tracingLog.String("id", id))

	var result entity.FlowSequence
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

func (repo *flowSequenceRepositoryImpl) Store(ctx context.Context, tenant string, entity *entity.FlowSequence) (*entity.FlowSequence, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceRepository.Store")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
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

func (repo flowSequenceRepositoryImpl) Delete(ctx context.Context, tenant, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceRepository.Delete")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("id", id))

	err := repo.gormDb.Delete(&entity.FlowSequence{}, "id = ?", id).Error

	if err != nil {
		span.LogFields(tracingLog.Bool("entity.deleted", false))
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(tracingLog.Bool("entity.deleted", true))

	return nil
}
