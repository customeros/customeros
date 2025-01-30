package events_listeners

import (
	"context"
	"fmt"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type MailstackProvisionMailboxListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewMailstackProvisionMailboxListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &MailstackProvisionBuyRequestListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.MailstackProvisionBuyRequest](), // subscribed event
			events.QueueEvents, // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *MailstackProvisionMailboxListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackProvisionMailboxListener.Handle")
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

	return l.handle(ctx, event.Event.EntityId)
}

func (l *MailstackProvisionMailboxListener) handle(ctx context.Context, entityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackProvisionMailboxListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	mailbox, err := l.dependencies.PostgresRepositories.TenantSettingsMailboxRepository.GetById(ctx, entityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = l.dependencies.CommonServices.OpenSRSService.SetupMailbox(ctx, mailbox.Tenant, mailbox.MailboxUsername, mailbox.MailboxPassword, strings.Split(mailbox.ForwardingTo, ","), mailbox.WebmailEnabled)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	mailbox.Status = postgres_entity.MailboxStatusProvisioned

	err = l.dependencies.PostgresRepositories.CommonRepository.UpdateProperty(ctx, mailbox.Tenant, postgres_entity.TenantSettingsMailbox{}, mailbox.ID, "Status", mailbox.Status)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return nil
}

func (l *MailstackProvisionMailboxListener) validateMessage(ctx context.Context, event *dto.Event) (*dto.MailstackProvisionMailbox, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackProvisionMailboxListener.validateMessage")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message, ok := event.Event.Data.(*dto.MailstackProvisionMailbox)
	if !ok {
		err := fmt.Errorf("expected MailstackProvisionMailbox, got %T", event.Event.Data)
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
