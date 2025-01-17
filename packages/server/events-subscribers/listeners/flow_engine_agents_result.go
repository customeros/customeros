package listeners

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/handlers"
)

func OnFlowAgentExecutionResultsEventCreated(ctx context.Context, s *service.CommonServices, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnFlowAgentExecutionResultsEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	flowAgentResultsEvent, err := getFlowAgentResultsEvent(input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowAgentResultsEvent.FlowAgentExecutionID == "" {
		err := errors.New("flowAgentExecutionResultsEvent is empty")
		tracing.TraceErr(span, err)
		return err
	}

	err = handlers.HandleAgentExecutionResults(ctx, s, flowAgentResultsEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func getFlowAgentResultsEvent(input any) (*dto.FlowAgentExecutionResultEvent, error) {
	message, ok := input.(*dto.Event)
	if !ok {
		return nil, fmt.Errorf("failed to cast to Event")
	}
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		return nil, err
	}

	flowAgentResultsEvent, ok := message.Event.Data.(*dto.FlowAgentExecutionResultEvent)
	if !ok {
		err := errors.New("event is not a flow action event")
		return nil, err
	}

	return flowAgentResultsEvent, nil
}
