package interfaces

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"

	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AgentService interface {
	CreateAgent(ctx context.Context, agentType enum.AgentType) (*postgresentity.Agent, error)
	UpdateAgent(ctx context.Context, agentId string, agentFields data_fields.AgentFields, capabilities *postgresentity.CapabilitiesConfig) (*postgresentity.Agent, error)
	GetAgentById(ctx context.Context, agentId string) (*postgresentity.Agent, error)
	GetAllAgentsByTenant(ctx context.Context) ([]*postgresentity.Agent, error)

	CreateAgentExecutionRecord(ctx context.Context, agent postgresentity.Agent, triggerEvent, traceId string) (string, error)
	SaveAgentExecutionCompleted(ctx context.Context, executionID string, goalAchieved bool) error
	SaveAgentExecutionError(ctx context.Context, executionID, errorMessage string) error
}
