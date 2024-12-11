package repository

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type FlowDeadListenerEventsRepository interface {
	Save(ctx context.Context, deadEvent entity.FlowDeadListenerEvents) (string, error)
}

type flowDeadListenerEventsRepository struct {
	gormDb *gorm.DB
}

func NewFlowDeadListenerEventsRepository(gormDb *gorm.DB) FlowDeadListenerEventsRepository {
	return &flowDeadListenerEventsRepository{gormDb: gormDb}
}

func (f *flowDeadListenerEventsRepository) Save(ctx context.Context, deadEvent entity.FlowDeadListenerEvents) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowDeadListenerEventsRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if deadEvent.Tenant == "" || deadEvent.ListenerEvent == "" || deadEvent.EventType == "" {
		span.LogFields(log.Object("deadEvent", deadEvent))
		err := errors.New("Tenant, ListenerEvent, or EventType missing")
		tracing.TraceErr(span, err)
		return "", err
	}

	err := f.gormDb.Save(&deadEvent).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return deadEvent.ID, nil
}
