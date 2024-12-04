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

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/model"
)

func OnWebhookEventCreated(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnWebhookEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	tenant, event, err := getWebhookEvent(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	c := model.EventContext{
		Context:  ctx,
		Span:     span,
		Services: services,
		Tenant:   tenant,
	}

	// determine event handler //
	switch eventData := (*event.Data).(type) {

	case *data_fields.MeetingSummaryEvent:
		c.SourceSystem = event.ExternalSystemId
		c.SourceEvent = event.Name
		handlers.HandleMeetingSummaryEvent(c, eventData)

	default:
		err := fmt.Errorf("Unsupported event %s", event.Name)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func getWebhookEvent(ctx context.Context, input any) (string, *dto.WebhookEvent[any], error) {
	message := input.(*dto.Event)
	tenant := message.Event.Tenant
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		return tenant, nil, err
	}

	return tenant, message.Event.Data.(*dto.WebhookEvent[any]), nil
}
