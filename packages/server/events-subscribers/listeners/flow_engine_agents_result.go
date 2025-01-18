package listeners

import (
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
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

	err = handleAgentExecutionResults(ctx, s, flowAgentResultsEvent)
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

func handleAgentExecutionResults(c context.Context, s *service.CommonServices, eventData *dto.FlowAgentExecutionResultEvent) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.HandleAgentExecutionResults")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant: eventData.Tenant,
	})

	// get existing flow execution record from db
	flowExecutionRecord, err := s.WorkflowService.GetFlowExecutionRecordById(ctx, eventData.FlowExecutionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// update flow execution record based on results
	switch eventData.Status {
	case "FAIL":
		flowExecutionRecord.Status = enum.FlowExecutionError.String()
		flowExecutionRecord.ErrorMessage = eventData.ErrorMessage

	case "PENDING":
		flowExecutionRecord.Status = enum.FlowExecutionRunning.String()

	case "SUCCESS":

		nextStepData, err := s.WorkflowService.GetNextStepInFlow(ctx, flowExecutionRecord.FlowID, nil)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		// if next step is END
		if nextStepData.ToNodeType == enum.NodeFlowEnd {
			flowExecutionRecord.CompletedAt = utils.NowPtr()
			flowExecutionRecord.CurrentStepNodeId = ""
			flowExecutionRecord.Status = enum.FlowExecutionCompleted.String()
			_, err := s.WorkflowService.SaveFlowExecutionRecord(ctx, *flowExecutionRecord)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
			return nil
		}

		flowExecutionRecord.CurrentStepNodeId = nextStepData.ToNodeID
		flowExecutionRecord.Status = enum.FlowExecutionRunning.String()
		// todo - fire next event & handle waits

	default:
		err = errors.New("FlowAgentExecutionStatus not valid")
		tracing.TraceErr(span, err)
		return err
	}

	_, err = s.WorkflowService.SaveFlowExecutionRecord(ctx, *flowExecutionRecord)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
