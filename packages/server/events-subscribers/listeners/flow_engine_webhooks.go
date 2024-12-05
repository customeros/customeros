package listeners

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/model"
)

func OnWebhookEventCreated(c context.Context, s *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(c, "Listeners.OnWebhookEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	ctx, eventName, webhookEvent, err := getWebhookEvent(ctx, input)
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

	case "MeetingSummaryEvent":
		eventData, ok := webhookEvent.Data.(*data_fields.MeetingSummaryEvent)
		if !ok {
			return fmt.Errorf("failed to cast to MeetingSummaryEvent, got type: %T", webhookEvent.Data)
		}
		return handlers.HandleMeetingSummaryEvent(ctx, s, eventName, eventData)

	default:
		err := fmt.Errorf("Unsupported event %s", webhookEvent.Name)
		tracing.TraceErr(span, err)
		return err
	}
}

func getWebhookEvent(ctx context.Context, input any) (context.Context, commonenum.FlowEvent, *dto.WebhookEvent, error) {
	message := input.(*dto.Event)
	tenant := message.Event.Tenant
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		return ctx, commonenum.NotSet, nil, err
	}

	event, ok := message.Event.Data.(*dto.WebhookEvent)
	if !ok {
		err := errors.New("event is not a webhook event")
		return ctx, commonenum.NotSet, nil, err
	}

	// update context with tenant, pass this where tenant is needed
	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant: tenant,
	})

	return ctx, event.Name, event, nil
}
