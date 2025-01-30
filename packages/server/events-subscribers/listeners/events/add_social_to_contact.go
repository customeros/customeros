package events_listeners

import (
	"context"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

// AddSocialToContactListener handles when socials are added to contacts
type AddSocialToContactListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewAddSocialToContactListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
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

	data, err := events.DecodeEventData[dto.AddSocialToContact](ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if data.Social == "" {
		err = errors.New("Social not set on event message")
		tracing.TraceErr(span, err)
		return err
	}

	socialUrl := data.Social
	contactId := event.Event.EntityId

	if strings.Contains(data.Social, "linkedin.com") {
		err := l.dependencies.CommonServices.EnrichmentService.EnrichContact(ctx, contactId, socialUrl)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (l *AddSocialToContactListener) validateMessage(ctx context.Context, event *dto.Event) (*dto.AddSocialToContact, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AddSocialToContactListener.validateMessage")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	// Decode the data into the specific event type
	eventDataObj := dto.AddSocialToContact{}
	if err := mapstructure.Decode(event.Event.Data, &eventDataObj); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Failed to decode data"))
		return nil, err
	}

	if eventDataObj.Social == "" {
		err := errors.New("Social not set on event message")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &eventDataObj, nil
}
