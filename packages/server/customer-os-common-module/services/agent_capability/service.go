package agent_capability

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type CapabilityExecutionHandler func(ctx context.Context, executionContainer *dto.CapabilityExecutionContainer) error

type agentCapabilityService struct {
	executionHandlers    map[enum.AgentCapabilityType]CapabilityExecutionHandler
	postgresRepositories *postgres_repository.Repositories
	enrichmentService    interfaces.EnrichmentService
}

func NewAgentCapabilityService(
	postgresRepositories *postgres_repository.Repositories,
	enrichmentService interfaces.EnrichmentService,
) (interfaces.AgentCapabilityService, error) {
	service := agentCapabilityService{
		postgresRepositories: postgresRepositories,
		enrichmentService:    enrichmentService,
	}

	err := service.InitCapabilityRegistry()
	if err != nil {
		return nil, err
	}

	// Register handlers
	service.executionHandlers[enum.CapabilityIdentifyWebVisitor] = service.handleIdentifyWebsiteVisitorExecution

	return &service, nil
}
