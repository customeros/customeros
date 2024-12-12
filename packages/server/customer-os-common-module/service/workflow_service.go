package service

import (
	"fmt"

	"github.com/matoous/go-nanoid/v2"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
)

type WorkflowService interface {
	SaveWorkflow(ctx context.Context, workflow *entity.Workflow) (string, error)
	GetWorkflowsByListenerEvent(ctx context.Context, listenerEvent enum.FlowListenerEvent) ([]entity.Workflow, error)
	GetFirstAction(ctx context.Context, workflow *entity.Workflow) (enum.FlowAction, string, error)
	GetWorkflowByFlowID(ctx context.Context, flowID string) (entity.Workflow, error)
	SaveFlowExecution(ctx context.Context, executionRecord postgresEntity.FlowExecution) error
	GetFlowExecutionRecord(ctx context.Context, excutionID, tenant string) (postgresEntity.FlowExecution, error)
	GetNextStepInFlow(ctx context.Context, flowID, lastNodeId string) (nodeType, nodeAction, nodeID, error)
	SaveFlowActionExecution(ctx context.Context, actionExecutionRecord postgresEntity.ActionExecution) error
	IsFlowActionValidTransitionFromListener(ctx context.Context, fromNode enum.FlowListenerEvent, toNode enum.FlowAction) (bool, error)
	ValidateListener(ctx context.Context, listenerEvent enum.FlowListenerEvent) (bool, error)
	GenerateID(ctx context.Context, entity string) (string, error)
}

type workflowService struct {
	services *Services
}

func NewWorkflowService(services *Services) WorkflowService {
	return &workflowService{
		services: services,
	}
}

// todo
func (w *workflowService) SaveWorkflow(ctx context.Context, workflow *entity.Workflow) (string, error) {
	// neo4j - implement create and update rules
	return "", nil
}

func (w *workflowService) GenerateID(ctx context.Context, entity string) (string, error) {
	id, err := fmt.Sprintf("%s-%s", entity, gonanoid.New())
	if err != nil {
		return "", err
	}
	return id, nil
}

type FlowListenerEventRecord struct {
	System      string `json:"system"`
	Event       string `json:"event"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (w *workflowService) ValidateListener(ctx context.Context, listenerEvent enum.FlowListenerEvent) (bool, error) {

	span, ctx := tracing.StartTracerSpan(ctx, "WorkflowService.ValidateListener")
	defer span.Finish()
	tracing.TagComponentRest(span)

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

func (w *workflowService) GetWorkflowByFlowID(ctx context.Context, flowID string) (entity.Workflow, error) {
	// neo4j
	var flow entity.Workflow
	return flow, nil
}

func (w *workflowService) GetFlowExecutionRecord(ctx context.Context, excutionID, tenant string) (postgresEntity.FlowExecution, error) {
	var executionRecord postgresEntity.FlowExecution
	return executionRecord, nil
}

func (w *workflowService) IsFlowActionValidTransitionFromListener(ctx context.Context, fromNode enum.FlowListenerEvent, toNode enum.FlowAction) (bool, error) {
	return false, nil
}

func (w *workflowService) GetWorkflowsByListenerEvent(ctx context.Context, listenerEvent enum.FlowListenerEvent) ([]entity.Workflow, error) {
	return []entity.Workflow{}, nil
}

func (w *workflowService) GetFirstAction(ctx context.Context, workflow *entity.Workflow) (enum.FlowAction, string, error) {
	// Lookup Workflow, determine what comes after listener
	// Call flow execution to ensure action hasn't already triggered
	var action enum.FlowAction
	return action, "", nil
}

func (w *workflowService) SaveFlowExecution(ctx context.Context, executionRecord postgresEntity.FlowExecution) error {
	return nil
}

func (w *workflowService) SaveFlowActionExecution(ctx context.Context, actionExecutionRecord postgresEntity.ActionExecution) error {
	return nil
}
