package listeners

import (
	"fmt"
	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"

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

	tenant, webhookEvent, err := getWebhookEvent(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if webhookEvent.Data == nil {
		err := errors.New("webhookEvent.Data is nil")
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
	switch webhookEvent.Name {

	case commonenum.EventFathomMeetingSummaryCreated:
		c.SourceSystem = webhookEvent.ExternalSystemId
		c.SourceEvent = webhookEvent.Name
		return handlers.HandleMeetingSummaryEvent(c, (*webhookEvent.Data).(*data_fields.MeetingSummaryEvent))

	default:
		err := fmt.Errorf("Unsupported event %s", webhookEvent.Name)
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
