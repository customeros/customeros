package events_listeners

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type RequestRefreshLastTouchpointListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewRequestRefreshLastTouchpointListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &HideContactListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.RequestRefreshLastTouchpoint](), // subscribed event
			events.QueueEvents, // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *RequestRefreshLastTouchpointListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestRefreshLastTouchpointListener.Handle")
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

	orgId := event.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, orgId)

	err = l.dependencies.CommonServices.OrganizationService.RefreshLastTouchpoint(ctx, orgId)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return nil
}

func (l *RequestRefreshLastTouchpointListener) validateMessage(ctx context.Context, event *dto.Event) (*dto.RequestRefreshLastTouchpoint, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestRefreshLastTouchpointListener.validateMessage")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message, ok := event.Event.Data.(*dto.RequestRefreshLastTouchpoint)
	if !ok {
		err := fmt.Errorf("expected RequestRefreshLastTouchpoint, got %T", event.Event.Data)
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
