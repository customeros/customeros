package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type ServiceLineItemEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func (h *ServiceLineItemEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemEventHandler.OnCreate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData event.ServiceLineItemCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	serviceLineItemId := aggregate.GetServiceLineItemObjectID(evt.GetAggregateID(), eventData.Tenant)
	err := h.repositories.ServiceLineItemRepository.CreateForContract(ctx, eventData.Tenant, serviceLineItemId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	return nil
}

func (h *ServiceLineItemEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemEventHandler.OnUpdate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData event.ServiceLineItemUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	serviceLineItemId := aggregate.GetServiceLineItemObjectID(evt.GetAggregateID(), eventData.Tenant)
	err := h.repositories.ServiceLineItemRepository.Update(ctx, eventData.Tenant, serviceLineItemId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	return nil
}
