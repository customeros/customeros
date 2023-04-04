package event_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type GraphPhoneNumberEventHandler struct {
	Repositories *repository.Repositories
}

func (e *GraphPhoneNumberEventHandler) OnPhoneNumberCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphPhoneNumberEventHandler.OnPhoneNumberCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.PhoneNumberCreatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	phoneNumberId := aggregate.GetPhoneNumberAggregateID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.PhoneNumberRepository.CreatePhoneNumber(ctx, phoneNumberId, eventData)

	return err
}

func (e *GraphPhoneNumberEventHandler) OnPhoneNumberUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphPhoneNumberEventHandler.OnPhoneNumberUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.PhoneNumberUpdatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	phoneNumberId := aggregate.GetPhoneNumberAggregateID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.PhoneNumberRepository.UpdatePhoneNumber(ctx, phoneNumberId, eventData)

	return err
}
