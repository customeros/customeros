package listeners

import (
	"fmt"
	"reflect"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func OnFlowAgentEventCreated(ctx context.Context, s *service.CommonServices, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnFlowAgentEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	flowAgentEvent, err := getFlowAgentEvent(input)
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
	return nil
}

func getFlowAgentEvent(input any) (*dto.FlowAgentEvent, error) {
	message, ok := input.(*dto.Event)
	if !ok {
		return nil, fmt.Errorf("failed to cast to Event")
	}
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		return nil, err
	}

	flowAgentEvent, ok := message.Event.Data.(*dto.FlowAgentEvent)
	if !ok {
		err := errors.New("event is not a flow action event")
		return nil, err
	}

	flowAgentData, ok := flowAgentEvent.Data.(map[string]interface{})
	if !ok {
		err := errors.New("event data is not a map")
		return nil, err
	}

	flowAgentDataPtr := reflect.New(eventDataTypes[flowAgentEvent.DataType]).Interface()
	err := utils.Decode(flowAgentData, flowAgentDataPtr)
	if err != nil {
		return nil, err
	}

	flowAgentEvent.Data = flowAgentDataPtr

	return flowAgentEvent, nil
}
