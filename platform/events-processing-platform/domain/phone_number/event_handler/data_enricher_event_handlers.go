package event_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type DataEnricherPhoneNumberEventHandler struct {
	commands *domain.CommandServices
}

func (e *DataEnricherPhoneNumberEventHandler) OnPhoneNumberCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "dataEnricherProjection.OnPhoneNumberCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.PhoneNumberCreatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	// FIXME alexb implement me
	// Find E164 and invoke a command

	return nil
}
