package service

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
)

type AgentCapabilitiesService interface {
	NotifySlackWebsiteVisit(ctx context.Context, event *data_fields.WebsiteVisitEvent) error
}

type agentCapabilitiesService struct {
	services *Services
}

func NewAgentCapabilitiesService(services *Services) AgentCapabilitiesService {
	return &agentCapabilitiesService{
		services: services,
	}
}
