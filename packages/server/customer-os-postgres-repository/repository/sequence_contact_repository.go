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

type SequenceContactRepository interface {
	Get(ctx context.Context, tenant, sequenceId string) ([]*entity.SequenceContact, error)
	GetById(ctx context.Context, tenant, id string) (*entity.SequenceContact, error)

	Store(ctx context.Context, tenant string, entity *entity.SequenceContact) (*entity.SequenceContact, error)
	Delete(ctx context.Context, tenant, id string) error
}

type sequenceContactRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewSequenceContactRepository(gormDb *gorm.DB) SequenceContactRepository {
	return &sequenceContactRepositoryImpl{gormDb: gormDb}
}

func (r sequenceContactRepositoryImpl) Get(ctx context.Context, tenant, sequenceId string) ([]*entity.SequenceContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceContactRepository.Get")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("sequenceId", sequenceId))

	var result []*entity.SequenceContact
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

func (r sequenceContactRepositoryImpl) GetById(ctx context.Context, tenant, id string) (*entity.SequenceContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceContactRepository.GetById")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	span.LogFields(tracingLog.String("id", id))

	var result entity.SequenceContact
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

func (repo *sequenceContactRepositoryImpl) Store(ctx context.Context, tenant string, entity *entity.SequenceContact) (*entity.SequenceContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceContactRepository.Store")
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

func (repo *sequenceContactRepositoryImpl) Delete(ctx context.Context, tenant, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceContactRepository.Delete")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("id", id))

	err := repo.gormDb.Delete(&entity.SequenceContact{}, "id = ?", id).Error

	if err != nil {
		span.LogFields(tracingLog.Bool("entity.deleted", false))
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(tracingLog.Bool("entity.deleted", true))

	return nil
}
