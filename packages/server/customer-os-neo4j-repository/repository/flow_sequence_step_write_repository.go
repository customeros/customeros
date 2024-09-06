package repository

//
//import (
//	"context"
//	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
//	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
//	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
//	"github.com/opentracing/opentracing-go"
//	tracingLog "github.com/opentracing/opentracing-go/log"
//	"gorm.io/gorm"
//)
//
//type FlowSequenceStepRepository interface {
//	GetList(ctx context.Context, sequenceId string) ([]*entity.FlowSequenceStep, error)
//	GetById(ctx context.Context, id string) (*entity.FlowSequenceStep, error)
//
//	Store(ctx context.Context, entity *entity.FlowSequenceStep) (*entity.FlowSequenceStep, error)
//}
//
//type flowSequenceStepRepositoryImpl struct {
//	gormDb *gorm.DB
//}
//
//func NewFlowSequenceStepRepository(gormDb *gorm.DB) FlowSequenceStepRepository {
//	return &flowSequenceStepRepositoryImpl{gormDb: gormDb}
//}
//
//func (r flowSequenceStepRepositoryImpl) GetList(ctx context.Context, sequenceId string) ([]*entity.FlowSequenceStep, error) {
//	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceStepRepository.GetList")
//	defer span.Finish()
//	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
//
//	span.LogFields(tracingLog.String("sequenceId", sequenceId))
//
//	var result []*entity.FlowSequenceStep
//	err := r.gormDb.
//		Where("tenant = ? and sequence_id = ?", common.GetTenantFromContext(ctx), sequenceId).
//		Find(&result).
//		Error
//
//	if err != nil {
//		tracing.TraceErr(span, err)
//		return nil, err
//	}
//
//	span.LogFields(tracingLog.Int("result.count", len(result)))
//
//	return result, nil
//}
//
//func (r flowSequenceStepRepositoryImpl) GetById(ctx context.Context, id string) (*entity.FlowSequenceStep, error) {
//	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceStepRepository.GetById")
//	defer span.Finish()
//	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
//
//	span.LogFields(tracingLog.String("id", id))
//
//	var result entity.FlowSequenceStep
//	err := r.gormDb.
//		Where("tenant = ? and id = ?", common.GetTenantFromContext(ctx), id).
//		First(&result).
//		Error
//
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			span.LogFields(tracingLog.Bool("result.found", false))
//			return nil, nil
//		}
//		tracing.TraceErr(span, err)
//		return nil, err
//	}
//
//	span.LogFields(tracingLog.Bool("result.found", true))
//
//	return &result, nil
//}
//
//func (repo *flowSequenceStepRepositoryImpl) Store(ctx context.Context, entity *entity.FlowSequenceStep) (*entity.FlowSequenceStep, error) {
//	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceStepRepository.Store")
//	defer span.Finish()
//	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
//
//	span.LogFields(tracingLog.Object("entity", entity))
//
//	err := repo.gormDb.Save(entity).Error
//
//	if err != nil {
//		tracing.TraceErr(span, err)
//		return nil, err
//	}
//
//	span.LogFields(tracingLog.String("entity.id", entity.ID))
//
//	return entity, nil
//}
