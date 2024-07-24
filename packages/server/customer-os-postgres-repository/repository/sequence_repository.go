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

type SequenceRepository interface {
	Get(ctx context.Context, tenant string) ([]*entity.Sequence, error)
	GetById(ctx context.Context, tenant, id string) (*entity.Sequence, error)

	Store(ctx context.Context, entity *entity.Sequence) (*entity.Sequence, error)
}

type sequenceRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewSequenceRepository(gormDb *gorm.DB) SequenceRepository {
	return &sequenceRepositoryImpl{gormDb: gormDb}
}

func (r sequenceRepositoryImpl) Get(ctx context.Context, tenant string) ([]*entity.Sequence, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceRepository.Get")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	var result []*entity.Sequence
	err := r.gormDb.
		Find(&result).
		Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(result)))

	return result, nil
}

func (r sequenceRepositoryImpl) GetById(ctx context.Context, tenant, id string) (*entity.Sequence, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceRepository.GetById")
	defer span.Finish()

	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagComponent, constants.ComponentPostgresRepository)

	span.LogFields(tracingLog.String("id", id))

	var result entity.Sequence
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

func (repo *sequenceRepositoryImpl) Store(ctx context.Context, entity *entity.Sequence) (*entity.Sequence, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "SequenceRepository.Store")
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
