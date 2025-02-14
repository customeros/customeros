package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AgentService interface {
	SetListeners(any)

	CreateAgent(ctx context.Context, agentType enum.AgentType) (*postgres_entity.Agent, error)
	UpdateAgent(ctx context.Context, agentId string, agentFields data_fields.AgentFields, capabilities []postgres_entity.Capability, listeners []postgres_entity.Listener) (*postgres_entity.Agent, error)
	DeleteAgent(ctx context.Context, agentId string) error
	GetAgentById(ctx context.Context, agentId string) (*postgres_entity.Agent, error)
	GetAllAgents(ctx context.Context) ([]*postgres_entity.Agent, error)

	CreateAgentExecutionRecord(ctx context.Context, agent postgres_entity.Agent, triggerEvent, traceId string) (string, error)
	GetAgentInfo(ctx context.Context) (*map[enum.AgentType]AgentInfo, error)
}

type AgentRegistry interface {
	SyncRegistry(ctx context.Context) error
}

type AgentInfo struct {
	Goal        string
	Description string
}
