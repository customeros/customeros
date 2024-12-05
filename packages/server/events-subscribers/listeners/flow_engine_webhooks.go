package listeners

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"reflect"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/handlers"
)

var eventDataTypes = map[string]reflect.Type{
	data_fields.MeetingSummaryEvent{}.Type(): reflect.TypeOf(data_fields.MeetingSummaryEvent{}),
	data_fields.ContactCreateEvent{}.Type():  reflect.TypeOf(data_fields.ContactCreateEvent{}),
	data_fields.MarkdownEventFields{}.Type(): reflect.TypeOf(data_fields.MarkdownEventFields{}),
}

func OnWebhookEventCreated(ctx context.Context, s *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnWebhookEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	eventName, webhookEvent, err := getWebhookEvent(input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if webhookEvent.Data == nil {
		err := errors.New("webhookEvent.Data is nil")
		tracing.TraceErr(span, err)
		return err
	}

	// determine event handler //
	switch webhookEvent.DataType {

	case data_fields.MeetingSummaryEvent{}.Type():
		eventData, ok := webhookEvent.Data.(*data_fields.MeetingSummaryEvent)
		if !ok {
			return fmt.Errorf("failed to cast to MeetingSummaryEvent, got type: %T", webhookEvent.Data)
		}
		err = handlers.HandleMeetingSummaryEvent(ctx, s, eventName, eventData)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

	default:
		err = fmt.Errorf("Unsupported event %s", webhookEvent.Name)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func getWebhookEvent(input any) (commonenum.FlowEvent, *dto.WebhookEvent, error) {
	message := input.(*dto.Event)
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		return commonenum.NotSet, nil, err
	}

	webhookEvent, ok := message.Event.Data.(*dto.WebhookEvent)
	if !ok {
		err := errors.New("event is not a webhook event")
		return commonenum.NotSet, nil, err
	}

	webhookData, ok := webhookEvent.Data.(map[string]interface{})
	if !ok {
		err := errors.New("event data is not a map")
		return commonenum.NotSet, nil, err
	}

	webhookDataPtr := reflect.New(eventDataTypes[webhookEvent.DataType]).Interface()
	err := utils.Decode(webhookData, webhookDataPtr)
	if err != nil {
		return commonenum.NotSet, nil, err
	}

	webhookEvent.Data = webhookDataPtr

	return webhookEvent.Name, webhookEvent, nil
}
