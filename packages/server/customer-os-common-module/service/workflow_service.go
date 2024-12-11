package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
)

type WorkflowService interface {
	SaveWorkflow(ctx context.Context, workflow *entity.Workflow) (string, error)
	GetWorkflowsByListenerEvent(ctx context.Context, listenerEvent enum.FlowListenerEvent) ([]entity.Workflow, error)
	GetFirstAction(ctx context.Context, workflow *entity.Workflow) (enum.FlowAction, string, error)
	SaveFlowExecution(ctx context.Context, executionRecord postgresEntity.FlowExecution) error
	SaveFlowActionExecution(ctx context.Context, actionExecutionRecord postgresEntity.ActionExecution) error
	IsFlowActionValidTransitionFromListener(ctx context.Context, fromNode enum.FlowListenerEvent, toNode enum.FlowAction) (bool, error)
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
	return "", nil
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
