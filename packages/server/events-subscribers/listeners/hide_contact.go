package listeners

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
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

func NewHideContactListener(logger logger.Logger, deps *model.DependencyContainer) events.EventListener {
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

	_, err = l.validateMessage(ctx, event)
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

func (l *HideContactListener) validateMessage(ctx context.Context, event *dto.Event) (*dto.HideContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HideContactListener.validateMessage")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message, ok := event.Event.Data.(*dto.HideContact)
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
