package service

import (
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

type WorkflowService interface {
	ValidateEventType(ctx context.Context, nodeType enum.FlowNodeType, event string) bool
	ValidateFlowBelongsToTenant(ctx context.Context, flowId string) bool
	ValidateListener(ctx context.Context, listenerEvent enum.FlowListenerEvent) (bool, error)
	ValidateNodeType(ctx context.Context, nodeType string) (bool, *enum.FlowNodeType)
	ValidateTransition(ctx context.Context, fromNodeId string, toNodeId string) (bool, error)
}

type workflowService struct {
	services *Services
}

func NewWorkflowService(services *Services) WorkflowService {
	return &workflowService{
		services: services,
	}
}
