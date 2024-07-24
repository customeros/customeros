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

type SequenceSenderRepository interface {
	Get(ctx context.Context, tenant, sequenceId string) ([]*entity.SequenceSender, error)
	GetById(ctx context.Context, tenant, id string) (*entity.SequenceSender, error)

	Store(ctx context.Context, tenant string, entity *entity.SequenceSender) (*entity.SequenceSender, error)
	Delete(ctx context.Context, tenant, id string) error
}

type sequenceSenderRepository struct {
	gormDb *gorm.DB
}

func NewSequenceSenderRepository(gormDb *gorm.DB) SequenceSenderRepository {
	return &sequenceSenderRepository{gormDb: gormDb}
}

func (r sequenceSenderRepository) Get(ctx context.Context, tenant, sequenceId string) ([]*entity.SequenceSender, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceSenderRepository.Get")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("sequenceId", sequenceId))

	var result []*entity.SequenceSender
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

func (r sequenceSenderRepository) GetById(ctx context.Context, tenant, id string) (*entity.SequenceSender, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceSenderRepository.GetById")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	span.LogFields(tracingLog.String("id", id))

	var result entity.SequenceSender
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

func (repo sequenceSenderRepository) Store(ctx context.Context, tenant string, entity *entity.SequenceSender) (*entity.SequenceSender, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceSenderRepository.Store")
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

func (repo sequenceSenderRepository) Delete(ctx context.Context, tenant, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceSenderRepository.Delete")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)
	span.LogFields(tracingLog.String("id", id))

	err := repo.gormDb.Delete(&entity.SequenceSender{}, "id = ?", id).Error

	if err != nil {
		span.LogFields(tracingLog.Bool("entity.deleted", false))
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(tracingLog.Bool("entity.deleted", true))

	return nil
}
