package listeners

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type RequestEnrichContactListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewRequestEnrichContactListener(logger logger.Logger, deps *model.DependencyContainer) events.EventListener {
	return &RequestEnrichContactListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.RequestEnrichContact](), // subscribed event
			events.QueueEvents, // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *RequestEnrichContactListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestEnrichContactListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = l.validateMessage(ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	contactId := event.Event.EntityId

	err = l.dependencies.CommonServices.EnrichmentService.EnrichContact(ctx, contactId, "")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (l *RequestEnrichContactListener) validateMessage(ctx context.Context, event *dto.Event) (*dto.RequestEnrichContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestEnrichContactListener.validateMessage")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message, ok := event.Event.Data.(*dto.RequestEnrichContact)
	if !ok {
		err := fmt.Errorf("expected RequestEnrichContact, got %T", event.Event.Data)
		tracing.TraceErr(span, err)
		return nil, err
	}

	if event.Event.EntityId == "" {
		err := errors.New("EntityId not set on event")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return message, nil
}
