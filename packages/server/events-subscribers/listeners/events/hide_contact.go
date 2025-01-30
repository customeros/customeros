package events_listeners

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type HideContactListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewHideContactListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &HideContactListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.HideContact](), // subscribed event
			events.QueueEvents,                     // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *HideContactListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HideContactListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	contactId := event.Event.EntityId

	err = l.dependencies.Neo4jRepositories.OrganizationWriteRepository.RefreshContactCountByContactId(ctx, nil, common.GetTenantFromContext(ctx), contactId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "OrganizationWriteRepository.RefreshContactCountByContactId"))
	}

	return nil
}
