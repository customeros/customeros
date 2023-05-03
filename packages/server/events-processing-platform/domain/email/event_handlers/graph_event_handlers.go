package event_handlers

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type GraphEmailEventHandler struct {
	Repositories *repository.Repositories
}

func (e *GraphEmailEventHandler) OnEmailCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphEmailEventHandler.OnEmailCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.EmailCreatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailId := aggregate.GetEmailAggregateID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.EmailRepository.CreateEmail(ctx, emailId, eventData)

	return err
}

func (e *GraphEmailEventHandler) OnEmailUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphEmailEventHandler.OnEmailUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.EmailUpdatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailId := aggregate.GetEmailAggregateID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.EmailRepository.UpdateEmail(ctx, emailId, eventData)

	return err
}

func (e *GraphEmailEventHandler) OnEmailValidationFailed(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphEmailEventHandler.OnEmailValidationFailed")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.EmailFailedValidationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailId := aggregate.GetEmailAggregateID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.EmailRepository.FailEmailValidation(ctx, emailId, eventData)

	return err
}

func (e *GraphEmailEventHandler) OnEmailValidated(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GraphEmailEventHandler.OnEmailValidated")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.EmailValidatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailId := aggregate.GetEmailAggregateID(evt.AggregateID, eventData.Tenant)
	err := e.Repositories.EmailRepository.ValidateEmail(ctx, emailId, eventData)

	return err
}
