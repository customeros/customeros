package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type EmailEventHandler struct {
	Repositories *repository.Repositories
}

func NewEmailEventHandler(repositories *repository.Repositories) *EmailEventHandler {
	return &EmailEventHandler{
		Repositories: repositories,
	}
}

func (h *EmailEventHandler) OnEmailCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.EmailCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailId := aggregate.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.EmailRepository.CreateEmail(ctx, emailId, eventData)

	return err
}

func (h *EmailEventHandler) OnEmailUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.EmailUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailId := aggregate.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.EmailRepository.UpdateEmail(ctx, emailId, eventData)

	return err
}

func (h *EmailEventHandler) OnEmailValidationFailed(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailValidationFailed")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.EmailFailedValidationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailId := aggregate.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.EmailRepository.FailEmailValidation(ctx, emailId, eventData)

	return err
}

func (h *EmailEventHandler) OnEmailValidated(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailValidated")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.EmailValidatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailId := aggregate.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	err := h.Repositories.EmailRepository.EmailValidated(ctx, emailId, eventData)

	return err
}
