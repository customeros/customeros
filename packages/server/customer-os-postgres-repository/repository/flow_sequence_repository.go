package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type FlowSequenceRepository interface {
	GetList(ctx context.Context, flowId *string) ([]*entity.FlowSequence, error)
	GetById(ctx context.Context, id string) (*entity.FlowSequence, error)
	GetFlowBySequenceId(ctx context.Context, id string) (*entity.Flow, error)

	Store(ctx context.Context, entity *entity.FlowSequence) (*entity.FlowSequence, error)
}

type flowSequenceRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewFlowSequenceRepository(gormDb *gorm.DB) FlowSequenceRepository {
	return &flowSequenceRepositoryImpl{gormDb: gormDb}
}

func (r flowSequenceRepositoryImpl) GetList(ctx context.Context, flowId *string) ([]*entity.FlowSequence, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceRepository.GetList")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.String("flowId", utils.StringOrEmpty(flowId)))

	var result []*entity.FlowSequence

	db := r.gormDb

	if flowId != nil {
		db = db.Where("tenant = ? and flow_id = ?", common.GetTenantFromContext(ctx), flowId)
	}

	err :=
		db.
			Find(&result).
			Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(result)))

	return result, nil
}

func (r flowSequenceRepositoryImpl) GetById(ctx context.Context, id string) (*entity.FlowSequence, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.String("id", id))

	var result entity.FlowSequence
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

func (r flowSequenceRepositoryImpl) GetFlowBySequenceId(ctx context.Context, id string) (*entity.Flow, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceRepository.GetFlowBySequenceId")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.String("id", id))

	var result entity.Flow
	err := r.gormDb.
		Table("flow_sequence").
		Select("flow.*").
		Joins("join flow on flow.id = flow_sequence.flow_id").
		Where("flow_sequence.tenant = ? and flow_sequence.id = ?", common.GetTenantFromContext(ctx), id).
		First(&result).
		Error
	if err != nil {
		return nil, err
	}

	if result.ID == "" {
		span.LogFields(tracingLog.Bool("result.found", false))
		return nil, nil
	}

	span.LogFields(tracingLog.Bool("result.found", true))

	return &result, nil
}

func (repo *flowSequenceRepositoryImpl) Store(ctx context.Context, entity *entity.FlowSequence) (*entity.FlowSequence, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceRepository.Store")
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
