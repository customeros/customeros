package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"golang.org/x/net/context"
)

type WorkflowService interface {
	ValidateEventType(ctx context.Context, nodeType enum.FlowNodeType, event string) bool
	ValidateFlowBelongsToTenant(ctx context.Context, flowId string) bool
	ValidateListener(ctx context.Context, listenerEvent enum.FlowListenerEvent) (bool, error)
	ValidateNodeType(ctx context.Context, nodeType string) (bool, *enum.FlowNodeType)
}

type workflowService struct {
	services *Services
}

func NewWorkflowService(services *Services) WorkflowService {
	return &workflowService{
		services: services,
	}
}
