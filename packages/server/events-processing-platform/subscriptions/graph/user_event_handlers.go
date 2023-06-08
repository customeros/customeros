package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type GraphUserEventHandler struct {
	Repositories *repository.Repositories
}

func (e *GraphUserEventHandler) OnUserCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphUserEventHandler.OnUserCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.UserCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.UserRepository.CreateUser(ctx, userId, eventData)

	return err
}

func (e *GraphUserEventHandler) OnUserUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphUserEventHandler.OnUserUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.UserUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.UserRepository.UpdateUser(ctx, userId, eventData)

	return err
}

func (e *GraphUserEventHandler) OnPhoneNumberLinkedToUser(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphUserEventHandler.OnPhoneNumberLinkedToUser")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.UserLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.PhoneNumberRepository.LinkWithUser(ctx, eventData.Tenant, userId, eventData.PhoneNumberId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (e *GraphUserEventHandler) OnEmailLinkedToUser(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphUserEventHandler.OnEmailLinkedToUser")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.UserLinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	userId := aggregate.GetUserObjectID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.EmailRepository.LinkWithUser(ctx, eventData.Tenant, userId, eventData.EmailId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}
