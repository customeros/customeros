package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func (a *agentService) TimelineAgent(ctx context.Context, event *dto.FlowAgentEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.TimelineAgent")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	switch event.DataType {
	case "MarkdownEventFields":
		eventData, ok := event.Data.(*data_fields.MarkdownEventFields)
		if !ok {
			return fmt.Errorf("failed to cast to MarkdownEventFields, got type: %T", event.Data)
		}
		if eventData == nil {
			return fmt.Errorf("MarkdownEventFields is nil")
		}

		return a.createMarkdownEvent(ctx, eventData, event.FlowExecutionId)

	default:
		return errors.New("Unsupported event")
	}
}

func (a *agentService) createMarkdownEvent(c context.Context, eventData *data_fields.MarkdownEventFields, flowExecutionID string) error {
	span, ctx := opentracing.StartSpanFromContext(c, "AgentService.createMarkdownEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	// create markdown event on timeline
	id, err := a.services.MarkdownEventService.Save(ctx, nil, nil, *eventData)

	// build execution record
	executionRecord := entity.FlowAgentExecution{
		FlowExecutionID: flowExecutionID,
		Agent:           enum.AgentTimelineEventCreate.String(),
		FlowNodeID:      "",
		StartedAt:       utils.NowPtr(),
	}

	if err != nil {
		executionRecord.Status = enum.FlowAgentExecutionFail.String()
		errMessage := fmt.Sprintf("Unable to save markdown event: %v", err)
		executionRecord.ErrorMessage = &errMessage
		tracing.TraceErr(span, err)

	} else {
		executionRecord.Status = enum.FlowAgentExecutionSuccess.String()
		executionRecord.CompletedAt = utils.NowPtr()
		result := fmt.Sprintf("Markdown Event ID: %s", id)
		executionRecord.Result = &result
	}

	// write action execution to db
	actionExecutionRecord, saveErr := a.services.PostgresRepositories.FlowAgentExecutionRepository.Create(ctx, executionRecord)
	if saveErr != nil {
		tracing.TraceErr(span, saveErr)
	}

	// fire action completion event
	pubErr := a.publishAgentResultEvent(ctx, flowExecutionID, actionExecutionRecord.ID, enum.FlowAgentExecutionStatus(executionRecord.Status), executionRecord.ErrorMessage)

	return multierr.Combine(err, saveErr, pubErr)
}
