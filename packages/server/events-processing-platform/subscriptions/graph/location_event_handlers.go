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

type LocationEventHandler struct {
	Repositories *repository.Repositories
}

func NewLocationEventHandler(repositories *repository.Repositories) *LocationEventHandler {
	return &LocationEventHandler{
		Repositories: repositories,
	}
}

func (h *LocationEventHandler) OnLocationCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationEventHandler.OnLocationCreate")
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

func (h *LocationEventHandler) OnLocationUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationEventHandler.OnLocationUpdate")
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

func (e *LocationEventHandler) OnLocationValidated(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationEventHandler.OnLocationValidated")
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

func (h *LocationEventHandler) OnLocationValidationFailed(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationEventHandler.OnLocationValidationFailed")
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
