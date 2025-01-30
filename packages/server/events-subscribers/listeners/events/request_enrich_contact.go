package events_listeners

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type RequestEnrichContactListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewRequestEnrichContactListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
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

	contactId := event.Event.EntityId
	span.LogKV("contactId", contactId)

	err = l.dependencies.CommonServices.EnrichmentService.EnrichContact(ctx, contactId, "")
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
