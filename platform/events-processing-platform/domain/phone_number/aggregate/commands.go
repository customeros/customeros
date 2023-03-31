package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *PhoneNumberAggregate) CreatePhoneNumber(ctx context.Context, tenant, rawPhoneNumber string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.CreatePhoneNumber")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	event, err := events.NewPhoneNumberCreatedEvent(a, tenant, rawPhoneNumber)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberCreatedEvent")
	}

	if err := event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
}
