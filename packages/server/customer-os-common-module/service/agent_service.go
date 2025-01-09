package service

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
)

type AgentService interface {
	VisitorIDAgent(ctx context.Context, eventData *data_fields.WebsiteVisitEvent)
	ICPAgent(ctx context.Context, event *dto.FlowAgentEvent) error
}

type agentService struct {
	services *Services
}

func NewAgentService(services *Services) AgentService {
	return &agentService{
		services: services,
	}
}
