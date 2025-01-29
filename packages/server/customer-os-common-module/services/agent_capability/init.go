package agent_capability

import (
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type AgentCapabilities struct {
	executors map[enum.AgentCapabilityType]interfaces.AgentCapabilityUntyped
}

func InitCapabilities(
	postgresRepositories *postgres_repository.Repositories,
	actionService interfaces.ActionService,
	aiService interfaces.AIService,
	enrichmentService interfaces.EnrichmentService,
	notificationService interfaces.NotificationService,
	organizationService interfaces.OrganizationService,
	tagService interfaces.TagService,
	workspaceService interfaces.WorkspaceService,

) *AgentCapabilities {

	capabilities := AgentCapabilities{}

	executors := make(map[enum.AgentCapabilityType]interfaces.AgentCapabilityUntyped)

	// Register each capability with its corresponding type.
	executors[enum.CapabilityAnalyzeWebSessionIntent] = NewAnalyzeWebSessionCapability(postgresRepositories, actionService)
	executors[enum.CapabilityCreateOrganization] = NewCreateOrganizationCapability(organizationService)
	executors[enum.CapabilityIdentifyWebVisitor] = NewIdentifyWebsiteVisitorCapability(postgresRepositories, enrichmentService)
	executors[enum.CapabilitySendSlackNotification] = NewSendSlackNotificationCapability(notificationService)
	executors[enum.CapabilitySendWebVisitorSlackNotification] = NewSendWebVisitorSlackNotificationCapability(postgresRepositories, notificationService, workspaceService)
	executors[enum.CapabilityApplyTag] = NewApplyTagCapability(tagService)
	// Continue registering other capabilities here...
	capabilities.executors = executors

	return &capabilities
}

// GetExecutor retrieves the untyped executor based on the capability type.
// Returns an error if the capability type is unsupported.
func (c *AgentCapabilities) GetExecutor(capType enum.AgentCapabilityType) (interfaces.AgentCapabilityUntyped, error) {
	executor, exists := c.executors[capType]
	if !exists {
		return nil, fmt.Errorf("unsupported capability type: %s", capType)
	}
	return executor, nil
}
