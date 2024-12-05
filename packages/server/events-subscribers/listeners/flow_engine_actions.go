package listeners

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/handlers"
)

func OnFlowActionEventCreated(c context.Context, s *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(c, "Listeners.OnFlowActionEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(c, span)
	tracing.LogObjectAsJson(span, "input", input)

	ctx, flowActionEvent, err := getFlowActionEvent(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowActionEvent.Data == nil {
		err := errors.New("flowActionEvent.Data is nil")
		tracing.TraceErr(span, err)
		return err
	}

	// determine event handler //
	switch flowActionEvent.DataType {

	case "MarkdownEventFields":
		eventData, ok := flowActionEvent.Data.(*data_fields.MarkdownEventFields)
		if !ok {
			return fmt.Errorf("failed to cast to MarkdownEventFields, got type: %T", flowActionEvent.Data)
		}
		return handlers.HandleCreateMarkdownEvent(c, s, eventData)

	case "ContactCreateEvent":
		eventData, ok := flowActionEvent.Data.(*data_fields.ContactCreateEvent)
		if !ok {
			return fmt.Errorf("failed to cast to ContactCreateEvent, got type: %T", flowActionEvent.Data)
		}
		return handlers.HandleCreateContact(c, s, eventData)

	default:
		err := fmt.Errorf("Unsupported flow action event %s", flowActionEvent.Name)
		tracing.TraceErr(span, err)
		return err
	}
}

func getFlowActionEvent(c context.Context, input any) (context.Context, *dto.FlowActionEvent, error) {
	message := input.(*dto.Event)
	tenant := message.Event.Tenant
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		return c, nil, err
	}

	// update context with tenant, pass this where tenant is needed
	ctx := common.WithCustomContext(c, &common.CustomContext{
		Tenant: tenant,
	})

	eventData, ok := message.Event.Data.(*dto.FlowActionEvent)
	if !ok {
		err := errors.New("event is not a webhook event")
		return ctx, nil, err
	}

	return ctx, eventData, nil
}
