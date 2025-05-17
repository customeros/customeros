package events_listeners

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type RequestRefreshLastTouchpointListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewRequestRefreshLastTouchpointListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &RequestRefreshLastTouchpointListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.RequestRefreshLastTouchpoint](), // subscribed event
			events.QueueEvents, // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *RequestRefreshLastTouchpointListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartListenerSpan(ctx, "RequestRefreshLastTouchpointListener.Handle")
	defer spans.Finish()
	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	orgId := event.Event.EntityId
	spans.TagEntity(orgId)

	err = l.dependencies.CommonServices.OrganizationService.RefreshLastTouchpoint(ctx, orgId)
	if err != nil {
		spans.TraceError(err)
	}

	return nil
}
