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

type SequenceStepRepository interface {
	Get(ctx context.Context, tenant, sequenceId string) ([]*entity.SequenceStep, error)
	GetById(ctx context.Context, tenant, id string) (*entity.SequenceStep, error)

	Store(ctx context.Context, tenant string, entity *entity.SequenceStep) (*entity.SequenceStep, error)
	Delete(ctx context.Context, tenant, id string) error
}

type sequenceStepRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewSequenceStepRepository(gormDb *gorm.DB) SequenceStepRepository {
	return &sequenceStepRepositoryImpl{gormDb: gormDb}
}

func (r sequenceStepRepositoryImpl) Get(ctx context.Context, tenant, sequenceId string) ([]*entity.SequenceStep, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceStepRepository.Get")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("sequenceId", sequenceId))

	var result []*entity.SequenceStep
	err := r.gormDb.
		Where("sequence_id = ?", sequenceId).
		Find(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(result)))

	return result, nil
}

func (r sequenceStepRepositoryImpl) GetById(ctx context.Context, tenant, id string) (*entity.SequenceStep, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceStepRepository.GetById")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	span.LogFields(tracingLog.String("id", id))

	var result entity.SequenceStep
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

func (repo *sequenceStepRepositoryImpl) Store(ctx context.Context, tenant string, entity *entity.SequenceStep) (*entity.SequenceStep, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceStepRepository.Store")
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

func (repo *sequenceStepRepositoryImpl) Delete(ctx context.Context, tenant, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceStepRepository.Delete")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("id", id))

	err := repo.gormDb.Delete(&entity.SequenceStep{}, "id = ?", id).Error

	if err != nil {
		span.LogFields(tracingLog.Bool("entity.deleted", false))
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(tracingLog.Bool("entity.deleted", true))

	return nil
}
