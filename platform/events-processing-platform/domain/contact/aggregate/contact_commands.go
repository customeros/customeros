package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/contact/events"
	"github.com/pkg/errors"
)

func (contactAggregate *ContactAggregate) CreateContact(ctx context.Context, uuid string, firstName string, lastName string) error {
	//span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.CreateContact")
	//defer span.Finish()
	//span.LogFields(log.String("AggregateID", contactAggregate.GetID()))

	if uuid == "" {
		return ErrorUuidIsRequired
	}

	contactCreatedEvent, err := events.NewContactCreatedEvent(contactAggregate, uuid, firstName, lastName)
	if err != nil {
		//tracing.TraceErr(span, err)
		//TODO VASI add logger

		return errors.Wrap(err, "NewOrderCreatedEvent")
	}

	//if err := contactCreatedEvent.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
	//	tracing.TraceErr(span, err)
	//	return errors.Wrap(err, "SetMetadata")
	//}

	return contactAggregate.Apply(contactCreatedEvent)
}

func (contactAggregate *ContactAggregate) UpdateContact(ctx context.Context, uuid string, firstName string, lastName string) error {
	//span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.UpdateContact")
	//defer span.Finish()
	//span.LogFields(log.String("AggregateID", contactAggregate.GetID()))

	orderUpdatedEvent, err := events.NewContactUpdatedEvent(contactAggregate, uuid, firstName, lastName)
	if err != nil {
		//tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactUpdatedEvent")
	}

	//if err := orderUpdatedEvent.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
	//	tracing.TraceErr(span, err)
	//	return errors.Wrap(err, "SetMetadata")
	//}

	return contactAggregate.Apply(orderUpdatedEvent)
}

func (contactAggregate *ContactAggregate) DeleteContact(ctx context.Context, uuid string) error {
	//span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.DeleteContact")
	//defer span.Finish()
	//span.LogFields(log.String("AggregateID", contactAggregate.GetID()))

	if uuid == "" {
		return ErrorUuidIsRequired
	}

	contactDeletedEvent, err := events.NewContactDeletedEvent(contactAggregate, uuid)
	if err != nil {
		//tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactDeletedEvent")
	}

	//if err := contactDeletedEvent.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
	//	tracing.TraceErr(span, err)
	//	return errors.Wrap(err, "SetMetadata")
	//}

	return contactAggregate.Apply(contactDeletedEvent)
}
