package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type GraphContactEventHandler struct {
	Repositories *repository.Repositories
}

func (e *GraphContactEventHandler) OnContactCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphContactEventHandler.OnContactCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.ContactCreatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.ContactRepository.CreateContact(ctx, contactId, eventData)

	return err
}

func (e *GraphContactEventHandler) OnContactUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphContactEventHandler.OnContactUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.ContactUpdatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.ContactRepository.UpdateContact(ctx, contactId, eventData)

	return err
}

func (e *GraphContactEventHandler) OnPhoneNumberLinkedToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphContactEventHandler.OnPhoneNumberLinkedToContact")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.ContactLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.PhoneNumberRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.PhoneNumberId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (e *GraphContactEventHandler) OnEmailLinkedToContact(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphContactEventHandler.OnEmailLinkedToContact")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.ContactLinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contactId := aggregate.GetContactObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.EmailRepository.LinkWithContact(ctx, eventData.Tenant, contactId, eventData.EmailId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}
