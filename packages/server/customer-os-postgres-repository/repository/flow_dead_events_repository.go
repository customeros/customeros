package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowDeadEventsRepository interface {
	Create(ctx context.Context, deadEvent entity.FlowDeadEvents) (*entity.FlowDeadEvents, error)
	Find(ctx context.Context, deadEvent entity.FlowDeadEvents) (*entity.FlowDeadEvents, error)
	Update(ctx context.Context, deadEvent entity.FlowDeadEvents) (*entity.FlowDeadEvents, error)
}

type flowDeadEventsRepository struct {
	gormDb *gorm.DB
}

func NewFlowDeadEventsRepository(gormDb *gorm.DB) FlowDeadEventsRepository {
	return &flowDeadEventsRepository{gormDb: gormDb}
}

func (f *flowDeadEventsRepository) Create(ctx context.Context, deadEvent entity.FlowDeadEvents) (*entity.FlowDeadEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowDeadEventsRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if deadEvent.Tenant == "" {
		deadEvent.Tenant = common.GetTenantFromContext(ctx)
		if deadEvent.Tenant == "" {
			err := errors.New("tenant missing from context")
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	if deadEvent.Event == "" || deadEvent.EventType == "" {
		span.LogFields(log.Object("deadEvent", deadEvent))
		err := errors.New("event, or event type missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	err := f.gormDb.Create(&deadEvent).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &deadEvent, nil
}

func (f *flowDeadEventsRepository) Find(ctx context.Context, deadEvent entity.FlowDeadEvents) (*entity.FlowDeadEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowDeadEventsRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if deadEvent.Tenant == "" {
		deadEvent.Tenant = common.GetTenantFromContext(ctx)
		if deadEvent.Tenant == "" {
			err := errors.New("tenant not set in context")
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	var foundDeadEvent entity.FlowDeadEvents
	err := f.gormDb.
		Where(&deadEvent).
		First(&foundDeadEvent).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &foundDeadEvent, nil
}

func (f *flowDeadEventsRepository) Update(ctx context.Context, deadEvent entity.FlowDeadEvents) (*entity.FlowDeadEvents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowDeadEventsRepository.Update")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if deadEvent.ID == "" {
		err := errors.New("ID is missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if deadEvent.Tenant == "" {
		deadEvent.Tenant = common.GetTenantFromContext(ctx)
		if deadEvent.Tenant == "" {
			err := errors.New("tenant not set in context")
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	var updatedDeadEvent entity.FlowDeadEvents
	err := f.gormDb.Model(&deadEvent).Updates(&deadEvent).First(&updatedDeadEvent).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &updatedDeadEvent, nil
}
