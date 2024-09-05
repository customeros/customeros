package repository

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type FlowSequenceContactRepository interface {
	GetList(ctx context.Context, sequenceId *string) ([]*entity.FlowSequenceContact, error)
	Identify(ctx context.Context, sequenceId, contactId, emailId string) (*entity.FlowSequenceContact, error)
	GetById(ctx context.Context, id string) (*entity.FlowSequenceContact, error)

	Store(ctx context.Context, entity *entity.FlowSequenceContact) (*entity.FlowSequenceContact, error)
	Delete(ctx context.Context, id string) error
}

type flowSequenceContactRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewFlowSequenceContactRepository(gormDb *gorm.DB) FlowSequenceContactRepository {
	return &flowSequenceContactRepositoryImpl{gormDb: gormDb}
}

func (r flowSequenceContactRepositoryImpl) GetList(ctx context.Context, sequenceId *string) ([]*entity.FlowSequenceContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactRepository.GetList")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	if sequenceId != nil {
		span.LogFields(tracingLog.String("sequenceId", *sequenceId))
	}

	var result []*entity.FlowSequenceContact

	db := r.gormDb

	if sequenceId != nil {
		db = db.Where("tenant = ? and sequence_id = ?", common.GetTenantFromContext(ctx), sequenceId)
	} else {
		db = db.Where("tenant = ?", common.GetTenantFromContext(ctx))
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

func (r flowSequenceContactRepositoryImpl) Identify(ctx context.Context, sequenceId, contactId, emailId string) (*entity.FlowSequenceContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactRepository.Identify")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.String("sequenceId", sequenceId), tracingLog.String("contactId", contactId), tracingLog.String("emailId", emailId))

	var result entity.FlowSequenceContact
	err := r.gormDb.
		Where("tenant = ? and sequence_id = ? and contact_id = ? and email_id = ?", common.GetTenantFromContext(ctx), sequenceId, contactId, emailId).
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

func (r flowSequenceContactRepositoryImpl) GetById(ctx context.Context, id string) (*entity.FlowSequenceContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.String("id", id))

	var result entity.FlowSequenceContact
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

func (repo *flowSequenceContactRepositoryImpl) Store(ctx context.Context, entity *entity.FlowSequenceContact) (*entity.FlowSequenceContact, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactRepository.Store")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.Object("entity", entity))

	if entity.Tenant == "" {
		tracing.TraceErr(span, errors.New("tenant is required"))
		return nil, errors.New("tenant is required")
	}

	err := repo.gormDb.Save(entity).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(tracingLog.String("entity.id", entity.ID))

	return entity, nil
}

func (repo *flowSequenceContactRepositoryImpl) Delete(ctx context.Context, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactRepository.Delete")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	span.LogFields(tracingLog.String("id", id))

	err := repo.gormDb.Delete(&entity.FlowSequenceContact{}, "tenant = ? and id = ?", common.GetTenantFromContext(ctx), id).Error

	if err != nil {
		span.LogFields(tracingLog.Bool("entity.deleted", false))
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(tracingLog.Bool("entity.deleted", true))

	return nil
}
