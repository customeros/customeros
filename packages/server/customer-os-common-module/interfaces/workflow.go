package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type WorkflowService interface {
	// Flow
	SaveFlow(ctx context.Context, flowRecord entity.Flows) (*entity.Flows, error)
	GetFlowsByTrigger(ctx context.Context, listenerEvent enum.FlowListenerEvent) ([]entity.Flows, error)
	GetNextStepInFlow(ctx context.Context, flowId string, fromNodeId *string) (*FlowNextStep, error)
	GetFlowsForListenerEvent(ctx context.Context, sourceEvent enum.FlowListenerEvent, eventType string, eventData any) ([]entity.Flows, error)

	// Flow Execution
	BuildAndSaveFlowExecutionRecord(ctx context.Context, flowStatus, flowId, flowNodeId, entityId, entityType, currentStep string, eventData any) (*entity.FlowExecution, error)
	SaveFlowExecutionRecord(ctx context.Context, flowExecutionRecord entity.FlowExecution) (*entity.FlowExecution, error)
	GetFlowExecutionRecordById(ctx context.Context, id string) (*entity.FlowExecution, error)

	// FlowAgent Execution
	SaveFlowAgentExecutionRecord(ctx context.Context, flowAgentExecutionRecord entity.AgentExecution) (*entity.AgentExecution, error)

	// Dead Flow Events
	// SendToDeadEvents(ctx context.Context, sourceEvent enum.FlowListenerEvent, eventType string, eventData any) error

	// Validation
	ValidateEventType(ctx context.Context, nodeType enum.FlowNodeType, event string) bool
	ValidateFlowBelongsToTenant(ctx context.Context, flowId string) (bool, error)
	ValidateListener(ctx context.Context, listenerEvent enum.FlowListenerEvent) (bool, error)
	ValidateNodeType(ctx context.Context, nodeType string) (bool, *enum.FlowNodeType)
	ValidateTransition(ctx context.Context, fromNodeId string, toNodeId string) (bool, error)
}

type FlowNextStep struct {
	FlowID      string
	FromNodeID  string
	ToNodeID    string
	ToNodeType  enum.FlowNodeType
	ToNodeAgent enum.FlowAgent
}
