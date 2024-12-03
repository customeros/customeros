package listeners

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func OnWebhookEventCreated(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnWebhookEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	event, err := getWebhookEvent(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// determine event handler //
	switch eventData := (*event.Data).(type) {
	case *data_fields.MeetingSummaryEvent:
		services.ActionMeetingSummaryService.HandleMeetingSummaryEvent(ctx, eventData, event.ExternalSystemId)
	default:
		err := fmt.Errorf("Unsupported event %s", event.Name)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func getWebhookEvent(ctx context.Context, input any) (*dto.WebhookEvent[any], error) {
	message := input.(*dto.Event)
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		return nil, err
	}

	return message.Event.Data.(*dto.WebhookEvent[any]), nil
}
