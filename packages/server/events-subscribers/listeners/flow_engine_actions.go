package listeners

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"reflect"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/handlers"
)

func OnFlowActionEventCreated(ctx context.Context, s *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnFlowActionEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	flowActionEvent, err := getFlowActionEvent(input)
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
		return handlers.HandleCreateMarkdownEvent(ctx, s, eventData)

	case "ContactCreateEvent":
		eventData, ok := flowActionEvent.Data.(*data_fields.ContactCreateEvent)
		if !ok {
			return fmt.Errorf("failed to cast to ContactCreateEvent, got type: %T", flowActionEvent.Data)
		}
		return handlers.HandleCreateContact(ctx, s, eventData)

	default:
		err := fmt.Errorf("Unsupported flow action event %s", flowActionEvent.Name)
		tracing.TraceErr(span, err)
		return err
	}
}

func getFlowActionEvent(input any) (*dto.FlowActionEvent, error) {
	message := input.(*dto.Event)
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		return nil, err
	}

	flowActionEvent, ok := message.Event.Data.(*dto.FlowActionEvent)
	if !ok {
		err := errors.New("event is not a webhook event")
		return nil, err
	}

	flowActionData, ok := flowActionEvent.Data.(map[string]interface{})
	if !ok {
		err := errors.New("event data is not a map")
		return nil, err
	}

	flowActionDataPtr := reflect.New(eventDataTypes[flowActionEvent.DataType]).Interface()
	err := utils.Decode(flowActionData, flowActionDataPtr)
	if err != nil {
		return nil, err
	}

	flowActionEvent.Data = flowActionDataPtr

	return flowActionEvent, nil
}
