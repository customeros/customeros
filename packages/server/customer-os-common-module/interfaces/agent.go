package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AgentService interface {
	CreateAgent(ctx context.Context, agentType enum.AgentType) (*postgres_entity.Agent, error)
	UpdateAgent(ctx context.Context, agentId string, agentFields data_fields.AgentFields, capabilities []postgres_entity.Capability, listeners []postgres_entity.Listener) (*postgres_entity.Agent, error)
	GetAgentById(ctx context.Context, agentId string) (*postgres_entity.Agent, error)
	GetAllAgentsByTenant(ctx context.Context) ([]*postgres_entity.Agent, error)

	CreateAgentExecutionRecord(ctx context.Context, agent postgres_entity.Agent, triggerEvent, traceId string) (string, error)
	SaveAgentExecutionCompleted(ctx context.Context, executionID string, goalAchieved bool) error
	SaveAgentExecutionError(ctx context.Context, executionID, errorMessage string) error
}

type AgentRegistry interface {
	SyncRegistry(ctx context.Context) error
}

type AgentEvents interface {
	Name() enum.AgentListenerEvent
}
