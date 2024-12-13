package handlers

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"
)

func HandleCreateMarkdownEvent(c context.Context, s *service.Services, eventData *data_fields.MarkdownEventFields, flowExecutionID string) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.HandleCreateMarkdownEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	// create markdown event on timeline
	id, err := s.MarkdownEventService.Save(ctx, nil, nil, *eventData)

	// build execution record
	executionRecord := entity.FlowActionExecution{
		FlowExecutionID: flowExecutionID,
		Action:          enum.ActionTimelineEventCreate.String(),
		FlowNodeID:      "",
		StartedAt:       utils.NowPtr(),
	}

	if err != nil {
		executionRecord.Status = enum.FlowActionExecutionFail.String()
		errMessage := fmt.Sprintf("Unable to save markdown event: %v", err)
		executionRecord.ErrorMessage = &errMessage
		tracing.TraceErr(span, err)

	} else {
		executionRecord.Status = enum.FlowActionExecutionSuccess.String()
		executionRecord.CompletedAt = utils.NowPtr()
		result := fmt.Sprintf("Markdown Event ID: %s", id)
		executionRecord.Result = &result
	}

	// write action execution to db
	actionExecutionRecord, saveErr := s.PostgresRepositories.FlowActionExecutionRepository.Create(ctx, executionRecord)
	if saveErr != nil {
		tracing.TraceErr(span, saveErr)
	}

	// fire action completion event
	pubErr := publishActionResultEvent(ctx, s, flowExecutionID, actionExecutionRecord.ID, enum.FlowActionExecutionStatus(executionRecord.Status), executionRecord.ErrorMessage)

	return multierr.Combine(err, saveErr, pubErr)
}

func publishActionResultEvent(
	ctx context.Context, s *service.Services, flowExecutionId, actionExecutionId string, actionExecutionStatus enum.FlowActionExecutionStatus, errorMessage *string,
) error {
	resultEvent := dto.FlowActionExecutionResultEvent{
		FlowExecutionID:       flowExecutionId,
		FlowActionExecutionID: actionExecutionId,
		Tenant:                common.GetTenantFromContext(ctx),
		Status:                actionExecutionStatus,
		ErrorMessage:          errorMessage,
	}

	s.RabbitMQService.PublishFlowActionEventResult(ctx, resultEvent)

	return nil
}
