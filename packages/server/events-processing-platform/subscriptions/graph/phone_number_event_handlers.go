package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type GraphPhoneNumberEventHandler struct {
	Repositories *repository.Repositories
}

func (h *GraphPhoneNumberEventHandler) OnPhoneNumberCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphPhoneNumberEventHandler.OnPhoneNumberCreate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.PhoneNumberCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.PhoneNumberRepository.CreatePhoneNumber(ctx, phoneNumberId, eventData)

	return err
}

func (h *GraphPhoneNumberEventHandler) OnPhoneNumberUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphPhoneNumberEventHandler.OnPhoneNumberUpdate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.PhoneNumberUpdatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.PhoneNumberRepository.UpdatePhoneNumber(ctx, phoneNumberId, eventData)

	return err
}

func (e *GraphPhoneNumberEventHandler) OnPhoneNumberValidated(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphPhoneNumberEventHandler.OnPhoneNumberValidated")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.PhoneNumberValidatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.PhoneNumberRepository.PhoneNumberValidated(ctx, phoneNumberId, eventData)

	return err
}

func (h *GraphPhoneNumberEventHandler) OnPhoneNumberValidationFailed(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphPhoneNumberEventHandler.OnPhoneNumberValidationFailed")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.PhoneNumberFailedValidationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	phoneNumberId := aggregate.GetPhoneNumberObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.PhoneNumberRepository.FailPhoneNumberValidation(ctx, phoneNumberId, eventData)

	return err
}
