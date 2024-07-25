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

type FlowSequenceStepRepository interface {
	Count(ctx context.Context, tenant, sequenceId string) (int64, error)
	Get(ctx context.Context, tenant, sequenceId string, page, limit int) ([]*entity.FlowSequenceStep, error)
	GetById(ctx context.Context, tenant, id string) (*entity.FlowSequenceStep, error)

	Store(ctx context.Context, tenant string, entity *entity.FlowSequenceStep) (*entity.FlowSequenceStep, error)
	Delete(ctx context.Context, tenant, id string) error
}

type flowSequenceStepRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewFlowSequenceStepRepository(gormDb *gorm.DB) FlowSequenceStepRepository {
	return &flowSequenceStepRepositoryImpl{gormDb: gormDb}
}

func (r flowSequenceStepRepositoryImpl) Count(ctx context.Context, tenant, sequenceId string) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceStepRepository.Count")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	var result int64
	err := r.gormDb.
		Model(entity.FlowSequenceStep{}).
		Where("sequence_id = ?", sequenceId).
		Count(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return 0, err
	}

	span.LogFields(tracingLog.Int64("result.count", result))

	return result, nil
}

func (r flowSequenceStepRepositoryImpl) Get(ctx context.Context, tenant, sequenceId string, page, limit int) ([]*entity.FlowSequenceStep, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceStepRepository.Get")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("sequenceId", sequenceId), tracingLog.Int("page", page), tracingLog.Int("limit", limit))

	var result []*entity.FlowSequenceStep
	err := r.gormDb.
		Where("sequence_id = ?", sequenceId).
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

func (r flowSequenceStepRepositoryImpl) GetById(ctx context.Context, tenant, id string) (*entity.FlowSequenceStep, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceStepRepository.GetById")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	span.LogFields(tracingLog.String("id", id))

	var result entity.FlowSequenceStep
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

func (repo *flowSequenceStepRepositoryImpl) Store(ctx context.Context, tenant string, entity *entity.FlowSequenceStep) (*entity.FlowSequenceStep, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceStepRepository.Store")
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

func (repo *flowSequenceStepRepositoryImpl) Delete(ctx context.Context, tenant, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceStepRepository.Delete")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("id", id))

	err := repo.gormDb.Delete(&entity.FlowSequenceStep{}, "id = ?", id).Error

	if err != nil {
		span.LogFields(tracingLog.Bool("entity.deleted", false))
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(tracingLog.Bool("entity.deleted", true))

	return nil
}
