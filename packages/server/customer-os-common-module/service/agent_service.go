package service

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

type AgentService interface {
	EmailAgent(ctx context.Context) error
	LinkedinAgent(ctx context.Context) error
	SlackAgent(ctx context.Context, event *dto.FlowAgentEvent) error
	TimelineAgent(ctx context.Context, event *dto.FlowAgentEvent) error
}

type agentService struct {
	services *Services
}

func NewAgentService(services *Services) AgentService {
	return &agentService{
		services: services,
	}
}

func (a *agentService) publishAgentResultEvent(
	ctx context.Context, flowExecutionId, actionExecutionId string, actionExecutionStatus enum.FlowAgentExecutionStatus, errorMessage *string,
) error {
	resultEvent := dto.FlowAgentExecutionResultEvent{
		FlowExecutionID:      flowExecutionId,
		FlowAgentExecutionID: actionExecutionId,
		Tenant:               common.GetTenantFromContext(ctx),
		Status:               actionExecutionStatus,
		ErrorMessage:         errorMessage,
	}

	a.services.RabbitMQService.PublishFlowAgentEventResult(ctx, resultEvent)

	return nil
}

func (a *agentService) EmailAgent(ctx context.Context) error {
	return nil
}

func (a *agentService) LinkedinAgent(ctx context.Context) error {
	return nil
}
