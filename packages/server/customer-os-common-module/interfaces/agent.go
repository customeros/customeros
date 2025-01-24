package interfaces

import (
	"context"

	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AgentService interface {
	CreateAgent(ctx context.Context, agentType enum.AgentType) (*postgresentity.Agents, error)
	CreateAgentExecutionRecord(ctx context.Context, agent postgresentity.Agents, triggerEvent string) (string, error)
	SaveAgentExecutionCompleted(ctx context.Context, executionID string, goalAchieved bool) error
	SaveAgentExecutionError(ctx context.Context, executionID, errorMessage string) error
}
