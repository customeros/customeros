package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

func (a *agentService) ICPAgent(ctx context.Context, event *dto.FlowAgentEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.ICPAgent")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	switch event.DataType {
	case "OrganizationQualifyEventFields":
		eventData, ok := event.Data.(*data_fields.OrganizationQualifyEventFields)
		if !ok {
			return fmt.Errorf("failed to cast to OrganizationQualifyEventFields, got type: %T", event.Data)
		}
		if eventData == nil {
			return fmt.Errorf("OrganizationQualifyEventFields is nil")
		}

		return a.buildICPQualificationReport(ctx, eventData, event.FlowExecutionId)

	default:
		return errors.New("Unsupported event")
	}
}

func (a *agentService) buildICPQualificationReport(ctx context.Context, eventData *data_fields.OrganizationQualifyEventFields, flowExecutionId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.buildICPQualificationReport")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	return nil
}
