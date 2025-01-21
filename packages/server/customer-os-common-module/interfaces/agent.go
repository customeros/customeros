package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AgentService interface {
	CreateAgent(ctx context.Context, agentType enum.AgentType) (*postgres_entity.Agents, error)
	CreateAgentExecutionRecord(ctx context.Context, agent postgres_entity.Agents, triggerEvent string) (string, error)
	SaveAgentExecutionSuccess(ctx context.Context, executionID string) error
	SaveAgentExecutionError(ctx context.Context, executionID, errorMessage string) error

	// DeleteAgent(ctx context.Context, agentID string) error
	// TurnOn(ctx context.Context, agentID string) error
	// TurnOff(ctx context.Context, agentID string) error
	// ListCapabilities(ctx context.Context, agentID string) ([]enum.AgentCapabilityType, error)
	// AddCapability(ctx context.Context, agentID string, capability enum.AgentCapabilityType)
	// RemoveCapability(ctx context.Context, agentID string, capability enum.AgentCapabilityType)
}
