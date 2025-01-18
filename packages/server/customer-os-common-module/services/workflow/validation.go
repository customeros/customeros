package workflow

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set on context")
		tracing.TraceErr(span, err)
		return false, err
	}

	return false, nil
}

//
// func (w *workflowService) ValidateListener(ctx context.Context, listenerEvent enum.FlowListenerEvent) (bool, error) {
// 	span, ctx := tracing.StartTracerSpan(ctx, "WorkflowService.ValidateListener")
// 	defer span.Finish()
// 	tracing.TagComponentService(span)
//
// 	tenant := common.GetTenantFromContext(ctx)
// 	if tenant == "" {
// 		err := errors.New("tenant not set on context")
// 		tracing.TraceErr(span, err)
// 		return false, err
// 	}
//
// 	events, err := w.services.PostgresRepositories.FlowListenerRegistryRepository.FindAll(ctx)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return false, err
// 	}
// 	for _, event := range *events {
// 		if event.ListenerEvent == listenerEvent.String() {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

func (w *workflowService) ValidateFlowBelongsToTenant(ctx context.Context, flowId string) (bool, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "WorkflowService.ValidateFlowBelongsToTenant")
	defer span.Finish()
	tracing.TagComponentService(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set on context")
		tracing.TraceErr(span, err)
		return false, err
	}

	query := entity.Flows{
		ID: flowId,
	}

	flowRecord, err := w.postgres.FlowsRepository.Find(ctx, query)
	if err == nil && flowRecord != nil {
		return true, nil
	}

	return false, nil
}

func (w *workflowService) ValidateEventType(ctx context.Context, nodeType enum.FlowNodeType, event string) bool {
	span, ctx := tracing.StartTracerSpan(ctx, "WorkflowService.ValidateEventType")
	defer span.Finish()
	tracing.TagComponentService(span)

	switch nodeType {
	case enum.NodeFlowEnd, enum.NodeFlowWait:
		return true
	case enum.NodeFlowAgent:
		_, err := enum.GetFlowAgent(event)
		if err == nil {
			return true
		}
	case enum.NodeFlowListenerEvent:
		_, err := enum.GetFlowListenerEvent(event)
		if err == nil {
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

func (w *workflowService) ValidateTransition(ctx context.Context, fromNodeId string, toNodeId string) (bool, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "WorkflowService.ValidateTransition")
	defer span.Finish()
	tracing.TagComponentService(span)

	fromNode, err := w.postgres.FlowNodeRepository.Find(ctx, entity.FlowNode{
		ID: fromNodeId,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	toNode, err := w.postgres.FlowNodeRepository.Find(ctx, entity.FlowNode{
		ID: toNodeId,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	results, err := w.postgres.FlowTransitionsRegistryRepository.Find(ctx, entity.FlowTransitionsRegistry{
		FromNodeType: fromNode.Type,
		ToNodeType:   toNode.Type,
	})

	if results == nil {
		err = errors.New("invalid transition")
		tracing.TraceErr(span, err)
		return false, err
	}
	return true, nil
}
