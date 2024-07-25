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

type FlowSequenceContactRepository interface {
	Get(ctx context.Context, tenant, sequenceId string) ([]*entity.FlowSequenceContact, error)
	GetById(ctx context.Context, tenant, id string) (*entity.FlowSequenceContact, error)

	Store(ctx context.Context, tenant string, entity *entity.FlowSequenceContact) (*entity.FlowSequenceContact, error)
	Delete(ctx context.Context, tenant, id string) error
}

type flowSequenceContactRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewFlowSequenceContactRepository(gormDb *gorm.DB) FlowSequenceContactRepository {
	return &flowSequenceContactRepositoryImpl{gormDb: gormDb}
}

func (r flowSequenceContactRepositoryImpl) Get(ctx context.Context, tenant, sequenceId string) ([]*entity.FlowSequenceContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactRepository.Get")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("sequenceId", sequenceId))

	var result []*entity.FlowSequenceContact
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

func (r flowSequenceContactRepositoryImpl) GetById(ctx context.Context, tenant, id string) (*entity.FlowSequenceContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactRepository.GetById")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	span.LogFields(tracingLog.String("id", id))

	var result entity.FlowSequenceContact
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

func (repo *flowSequenceContactRepositoryImpl) Store(ctx context.Context, tenant string, entity *entity.FlowSequenceContact) (*entity.FlowSequenceContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactRepository.Store")
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

func (repo *flowSequenceContactRepositoryImpl) Delete(ctx context.Context, tenant, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactRepository.Delete")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("id", id))

	err := repo.gormDb.Delete(&entity.FlowSequenceContact{}, "id = ?", id).Error

	if err != nil {
		span.LogFields(tracingLog.Bool("entity.deleted", false))
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(tracingLog.Bool("entity.deleted", true))

	return nil
}
