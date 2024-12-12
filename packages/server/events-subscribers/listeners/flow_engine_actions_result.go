package listeners

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/handlers"
)

func OnFlowActionExecutionResultsEventCreated(ctx context.Context, s *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnFlowActionExecutionResultsEventCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	flowActionResultsEvent, err := getFlowActionResultsEvent(input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if flowActionResultsEvent.FlowActionExecutionID == "" {
		err := errors.New("flowActionExecutionResultsEvent is empty")
		tracing.TraceErr(span, err)
		return err
	}

	err = handlers.HandleActionExecutionResults(ctx, s, flowActionResultsEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func getFlowActionResultsEvent(input any) (*dto.FlowActionExecutionResultEvent, error) {
	message, ok := input.(*dto.Event)
	if !ok {
		return nil, fmt.Errorf("failed to cast to Event")
	}
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		return nil, err
	}

	flowActionResultsEvent, ok := message.Event.Data.(*dto.FlowActionExecutionResultEvent)
	if !ok {
		err := errors.New("event is not a flow action event")
		return nil, err
	}

	return flowActionResultsEvent, nil
}
