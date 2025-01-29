package listeners

import (
	"context"
	"fmt"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

// AddSocialToContactListener handles when socials are added to contacts
type AddSocialToContactListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewAddSocialToContactListener(logger logger.Logger, deps *model.DependencyContainer) events.EventListener {
	return &AddSocialToContactListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.AddSocialToContact](), // subscribed event
			events.QueueEvents,                            // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *AddSocialToContactListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AddSocialToContactListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	message, err := l.validateMessage(ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	socialUrl := message.Social
	contactId := event.Event.EntityId

	if strings.Contains(message.Social, "linkedin.com") {
		err := l.dependencies.CommonServices.EnrichmentService.EnrichContact(ctx, contactId, socialUrl)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (l *AddSocialToContactListener) validateMessage(ctx context.Context, event *dto.Event) (*dto.AddSocialToContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AddSocialToContactListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message, ok := event.Event.Data.(*dto.AddSocialToContact)
	if !ok {
		err := fmt.Errorf("expected AddSocialToContact, got %T", event.Event.Data)
		tracing.TraceErr(span, err)
		return nil, err
	}

	if message.Social == "" {
		err := errors.New("Social not set on event message")
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
