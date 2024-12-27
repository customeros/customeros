package listeners

import (
	"fmt"
	"reflect"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func OnFlowAgentEventCreated(ctx context.Context, s *service.Services, input any) error {
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

	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant: flowAgentEvent.Tenant,
	})

	// determine event handler
	switch flowAgentEvent.Name.Agent() {

	case "slack":
		return s.AgentService.SlackAgent(ctx, flowAgentEvent)

	case "timeline_event":
		return s.AgentService.TimelineAgent(ctx, flowAgentEvent)

	default:
		err := fmt.Errorf("Unsupported flow action event %s", flowAgentEvent.Name)
		tracing.TraceErr(span, err)
		return err
	}
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
