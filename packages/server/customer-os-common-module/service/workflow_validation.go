package service

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type FlowListenerEventRecord struct {
	System      string `json:"system"`
	Event       string `json:"event"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (w *workflowService) ValidateListener(ctx context.Context, listenerEvent enum.FlowListenerEvent) (bool, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "WorkflowService.ValidateListener")
	defer span.Finish()
	tracing.TagComponentService(span)
	events, err := w.services.PostgresRepositories.FlowListenerRegistryRepository.GetAllFlowListenerEvents(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}
	for _, event := range events {
		if event.ListenerEvent == listenerEvent.String() {
			return true, nil
		}
	}
	return false, nil
}

func (w *workflowService) ValidateFlowBelongsToTenant(ctx context.Context, flowId string) bool {
	span, ctx := tracing.StartTracerSpan(ctx, "WorkflowService.ValidateFlowBelongsToTenant")
	defer span.Finish()
	tracing.TagComponentService(span)

	flowRecord, err := w.services.PostgresRepositories.FlowRepository.GetFlowByID(ctx, flowId)
	if err == nil && flowRecord != nil {
		return true
	}

	return false
}

func (w *workflowService) ValidateEventType(ctx context.Context, nodeType enum.FlowNodeType, event string) bool {
	span, ctx := tracing.StartTracerSpan(ctx, "WorkflowService.ValidateEventType")
	defer span.Finish()
	tracing.TagComponentService(span)

	switch nodeType {
	case enum.NodeFlowEnd, enum.NodeFlowWait:
		return true
	case enum.NodeFlowAction:
		_, err := enum.GetFlowAction(event)
		if err != nil {
			return true
		}
	case enum.NodeFlowListenerEvent:
		_, err := enum.GetFlowListenerEvent(event)
		if err != nil {
			return true
		}
	default:
		return false
	}

	return false

}

func (w *workflowService) ValidateNodeType(ctx context.Context, nodeType string) (bool, *enum.FlowNodeType) {
	span, ctx := tracing.StartTracerSpan(ctx, "WorkflowService.ValidateNodeType")
	defer span.Finish()
	tracing.TagComponentService(span)

	t, err := enum.GetFlowNodeType(nodeType)
	if err != nil {
		return false, nil
	}
	return true, &t
}
