package agent_capability

import (
	"context"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type CapabilityExecutionHandler func(ctx context.Context, executionContainer *dto.CapabilityExecutionContainer) error

type agentCapabilityService struct {
	executionHandlers    map[enum.AgentCapabilityType]CapabilityExecutionHandler
	postgresRepositories *postgres_repository.Repositories
	enrichmentService    interfaces.EnrichmentService
	organizationService  interfaces.OrganizationService
	actionService        interfaces.ActionService
	notificationService  interfaces.NotificationService
}

func NewAgentCapabilityService(
	postgresRepositories *postgres_repository.Repositories,
	enrichmentService interfaces.EnrichmentService,
	organizationService interfaces.OrganizationService,
	actionService interfaces.ActionService,
	notificationService interfaces.NotificationService,
) (interfaces.AgentCapabilityService, error) {
	service := &agentCapabilityService{
		executionHandlers:    make(map[enum.AgentCapabilityType]CapabilityExecutionHandler),
		postgresRepositories: postgresRepositories,
		enrichmentService:    enrichmentService,
		organizationService:  organizationService,
		actionService:        actionService,
		notificationService:  notificationService,
	}

	// err := service.InitCapabilityRegistry()
	// if err != nil {
	// 	return nil, err
	// }

	// Register handlers
	service.executionHandlers[enum.CapabilityAnalyzeWebSessionIntent] = service.handleAnalyzeWebSessionExecution
	service.executionHandlers[enum.CapabilityCreateOrganization] = service.handleOrganizationCreationExecution
	service.executionHandlers[enum.CapabilityIdentifyWebVisitor] = service.handleIdentifyWebsiteVisitorExecution

	return service, nil
}

func (c *agentCapabilityService) SetOrganizationService(org interfaces.OrganizationService) {
	c.organizationService = org
}

func (c *agentCapabilityService) SetActionService(action interfaces.ActionService) {
	c.actionService = action
}
