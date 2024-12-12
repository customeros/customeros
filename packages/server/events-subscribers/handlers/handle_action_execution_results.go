package handlers

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
)

func HandleActionExecutionResults(c context.Context, s *service.Services, eventData *dto.FlowActionExecutionResultEvent) error {
	span, ctx := opentracing.StartSpanFromContext(c, "EventHandlers.HandleCreateContact")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "eventData", eventData)

	// build Flow Execution update
	flowExecutionUpdate := entity.FlowExecution{
		ID:     eventData.FlowActionExecutionID,
		Tenant: eventData.Tenant,
	}

	if eventData.Status == enum.FlowActionExecutionFail {
		flowExecutionUpdate.Status = enum.FlowExecutionError.String()
		flowExecutionUpdate.ErrorMessage = eventData.ErrorMessage
		// don't change next step
	}

	// look for next action in flow, if none found mark Flow Execution as completed

	return nil
}
