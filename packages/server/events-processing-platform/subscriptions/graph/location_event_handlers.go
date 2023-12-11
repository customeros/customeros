package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type GraphLocationEventHandler struct {
	Repositories *repository.Repositories
}

func (h *GraphLocationEventHandler) OnLocationCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphLocationEventHandler.OnLocationCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.LocationCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	LocationId := aggregate.GetLocationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.LocationRepository.CreateLocation(ctx, LocationId, eventData)

	return err
}

func (h *GraphLocationEventHandler) OnLocationUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphLocationEventHandler.OnLocationUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.LocationUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	LocationId := aggregate.GetLocationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.LocationRepository.UpdateLocation(ctx, LocationId, eventData)

	return err
}

func (e *GraphLocationEventHandler) OnLocationValidated(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphLocationEventHandler.OnLocationValidated")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.LocationValidatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	LocationId := aggregate.GetLocationObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.LocationRepository.LocationValidated(ctx, LocationId, eventData)

	return err
}

func (h *GraphLocationEventHandler) OnLocationValidationFailed(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphLocationEventHandler.OnLocationValidationFailed")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.LocationFailedValidationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	LocationId := aggregate.GetLocationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.LocationRepository.FailLocationValidation(ctx, LocationId, eventData)

	return err
}
