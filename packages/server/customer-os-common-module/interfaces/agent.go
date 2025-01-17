package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type AgentService interface {
	CreateAgent(ctx context.Context) (*entity.Agents, error)
	RunAgent(ctx context.Context, agent *entity.Agents, eventData any) error
	AgentCapabilities(ctx context.Context, agent *entity.Agents) (*AgentCapabilities, error)

	SetActionService(action ActionService)
	SetOrganizationService(org OrganizationService)
	IsInitialized() bool
}

type AgentCapabilities struct {
	Capabilities []Capability `json:"capabilities"`
}

type Capability struct {
	Name     string `json:"name"`
	Action   string `json:"action"`
	Optional bool   `json:"optional,omitempty"`
}
