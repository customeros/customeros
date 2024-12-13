package handlers

import (
	"context"
	"errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
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
	case enum.FlowActionExecutionFail:
		flowExecutionRecord.Status = enum.FlowExecutionError.String()
		flowExecutionRecord.ErrorMessage = eventData.ErrorMessage

	case enum.FlowActionExecutionPending:
		flowExecutionRecord.Status = enum.FlowExecutionRunning.String()

	case enum.FlowActionExecutionSuccess:
		nextStepData, err := s.WorkflowService.GetNextStepInFlow(ctx, flowExecutionRecord.FlowID, flowExecutionRecord.CurrentStepNodeId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		// if next step is END
		if nextStepData.ToNodeType == enum.NodeFlowEnd {
			flowExecutionRecord.CompletedAt = utils.NowPtr()
			flowExecutionRecord.CurrentStep = ""
			flowExecutionRecord.CurrentStepNodeId = ""
			flowExecutionRecord.Status = enum.FlowExecutionCompleted.String()
			_, err := s.WorkflowService.SaveFlowExecutionRecord(ctx, *flowExecutionRecord)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
			return nil
		}

		flowExecutionRecord.CurrentStep = nextStepData.ToNodeAction.String()
		flowExecutionRecord.CurrentStepNodeId = nextStepData.ToNodeID
		flowExecutionRecord.Status = enum.FlowExecutionRunning.String()
		// todo - fire next event & handle waits

	default:
		err = errors.New("FlowActionExecutionStatus not valid")
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
