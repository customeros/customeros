package interfaces

import (
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
)

type AgentRunnerService interface {
	Run(ctx context.Context, agent postgres_entity.Agent, agentEventName string, initialParams map[string]any) error
}
