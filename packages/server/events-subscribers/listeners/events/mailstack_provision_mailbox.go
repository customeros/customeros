package events_listeners

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type MailstackProvisionMailboxListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewMailstackProvisionMailboxListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &MailstackProvisionMailboxListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.MailstackProvisionMailbox](), // subscribed event
			events.QueueEvents, // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *MailstackProvisionMailboxListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartListenerSpan(ctx, "MailstackProvisionMailboxListener.Handle")
	defer spans.Finish()
	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return l.handle(ctx, event.Event.EntityId)
}

func (l *MailstackProvisionMailboxListener) handle(ctx context.Context, entityId string) error {
	spans, ctx := telemetry.StartListenerSpan(ctx, "MailstackProvisionMailboxListener.handle")
	defer spans.Finish()
	spans.LogObjectAsJson("entityId", entityId)

	tenant := common.GetTenantFromContext(ctx)

	err := l.dependencies.CommonServices.MailstackService.ConfigureMailbox(ctx, tenant, entityId)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
