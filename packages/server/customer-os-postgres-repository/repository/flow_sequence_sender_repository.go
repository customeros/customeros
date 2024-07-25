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

type FlowSequenceSenderRepository interface {
	Count(ctx context.Context, tenant, sequenceId string) (int64, error)
	Get(ctx context.Context, tenant, sequenceId string, page, limit int) ([]*entity.FlowSequenceSender, error)
	GetById(ctx context.Context, tenant, id string) (*entity.FlowSequenceSender, error)

	Store(ctx context.Context, tenant string, entity *entity.FlowSequenceSender) (*entity.FlowSequenceSender, error)
	Delete(ctx context.Context, tenant, id string) error
}

type flowSequenceSenderRepository struct {
	gormDb *gorm.DB
}

func NewFlowSequenceSenderRepository(gormDb *gorm.DB) FlowSequenceSenderRepository {
	return &flowSequenceSenderRepository{gormDb: gormDb}
}

func (r flowSequenceSenderRepository) Count(ctx context.Context, tenant, sequenceId string) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceSenderRepository.Count")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	var result int64
	err := r.gormDb.
		Model(entity.FlowSequenceSender{}).
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

func (r flowSequenceSenderRepository) Get(ctx context.Context, tenant, sequenceId string, page, limit int) ([]*entity.FlowSequenceSender, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceSenderRepository.Get")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("sequenceId", sequenceId))

	var result []*entity.FlowSequenceSender
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

func (r flowSequenceSenderRepository) GetById(ctx context.Context, tenant, id string) (*entity.FlowSequenceSender, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceSenderRepository.GetById")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	span.LogFields(tracingLog.String("id", id))

	var result entity.FlowSequenceSender
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

func (repo flowSequenceSenderRepository) Store(ctx context.Context, tenant string, entity *entity.FlowSequenceSender) (*entity.FlowSequenceSender, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceSenderRepository.Store")
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

func (repo flowSequenceSenderRepository) Delete(ctx context.Context, tenant, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceSenderRepository.Delete")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("id", id))

	err := repo.gormDb.Delete(&entity.FlowSequenceSender{}, "id = ?", id).Error

	if err != nil {
		span.LogFields(tracingLog.Bool("entity.deleted", false))
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(tracingLog.Bool("entity.deleted", true))

	return nil
}
