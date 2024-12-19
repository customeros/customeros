package listeners

import (
	"fmt"
	"reflect"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/handlers"
)

func OnFlowAgentEventCreated(ctx context.Context, s *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnFlowAgentEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	flowAgentEvent, flowExecutionID, err := getFlowAgentEvent(input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowAgentEvent.Data == nil {
		err := errors.New("flowAgentEvent.Data is nil")
		tracing.TraceErr(span, err)
		return err
	}

	// determine event handler
	switch flowAgentEvent.DataType {

	case "MarkdownEventFields":
		eventData, ok := flowAgentEvent.Data.(*data_fields.MarkdownEventFields)
		if !ok {
			return fmt.Errorf("failed to cast to MarkdownEventFields, got type: %T", flowAgentEvent.Data)
		}
		return handlers.HandleCreateMarkdownEvent(ctx, s, eventData, flowExecutionID)

	case "ContactCreateEvent":
		eventData, ok := flowAgentEvent.Data.(*data_fields.ContactCreateEvent)
		if !ok {
			return fmt.Errorf("failed to cast to ContactCreateEvent, got type: %T", flowAgentEvent.Data)
		}
		return handlers.HandleCreateContact(ctx, s, eventData)

	default:
		err := fmt.Errorf("Unsupported flow action event %s", flowAgentEvent.Name)
		tracing.TraceErr(span, err)
		return err
	}
}

func getFlowAgentEvent(input any) (*dto.FlowAgentEvent, string, error) {
	message, ok := input.(*dto.Event)
	if !ok {
		return nil, "", fmt.Errorf("failed to cast to Event")
	}
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		return nil, "", err
	}

	flowAgentEvent, ok := message.Event.Data.(*dto.FlowAgentEvent)
	if !ok {
		err := errors.New("event is not a flow action event")
		return nil, "", err
	}

	flowExecutionID := flowAgentEvent.FlowExecutionId

	flowAgentData, ok := flowAgentEvent.Data.(map[string]interface{})
	if !ok {
		err := errors.New("event data is not a map")
		return nil, flowExecutionID, err
	}

	flowAgentDataPtr := reflect.New(eventDataTypes[flowAgentEvent.DataType]).Interface()
	err := utils.Decode(flowAgentData, flowAgentDataPtr)
	if err != nil {
		return nil, flowExecutionID, err
	}

	flowAgentEvent.Data = flowAgentDataPtr

	return flowAgentEvent, flowExecutionID, nil
}
