package handlers

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

func HandleActionExecutionResults(c context.Context, s *service.Services, eventData *dto.FlowActionExecutionResultEvent) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.HandleActionExecutionResults")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	// get existing flow execution record from db
	flowExecutionRecord, err := s.WorkflowService.GetFlowExecutionRecord(ctx, eventData.FlowExecutionID, eventData.Tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// update flow execution record based on results
	switch eventData.Status {
	case enum.FlowActionExecutionFail:
		flowExecutionRecord.Status = enum.FlowExecutionError.String()
		flowExecutionRecord.ErrorMessage = eventData.ErrorMessage

	case enum.FlowActionExecutionPending:
		flowExecutionRecord.Status = enum.FlowExecutionRunning.String()

	case enum.FlowActionExecutionSuccess:
		nextNodeType, nextNodeAction, nextNodeID, err := s.WorkflowService.GetNextStepInFlow(ctx, flowExecutionRecord.FlowID, flowExecutionRecord.CurrentStepNodeId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if nextNodeType == enum.NodeFlowEnd.String() {
			flowExecutionRecord.CompletedAt = utils.NowPtr()
			flowExecutionRecord.CurrentStep = ""
			flowExecutionRecord.CurrentStepNodeId = ""
			flowExecutionRecord.Status = enum.FlowExecutionCompleted.String()
		}

		//to do, handle next step
		flowExecutionRecord.CurrentStep = nextNodeAction
		flowExecutionRecord.CurrentStepNodeId = nextNodeID
		flowExecutionRecord.Status = enum.FlowExecutionRunning.String()
		// fire next event

	default:
		err = errors.New("FlowActionExecutionStatus not valid")
		tracing.TraceErr(span, err)
		return err
	}

	err = s.WorkflowService.SaveFlowExecution(ctx, flowExecutionRecord)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
