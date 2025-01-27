package agent_capability

import (
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type AgentCapabilities struct {
	AnalyzeWebSession      *AnalyzeWebSessionCapability
	ApplyTag               *ApplyTagCapability
	CreateOrganization     *CreateOrganizationCapability
	ICPQualification       *ICPQualificationCapability
	IdentifyWebsiteVisitor *IdentifyWebsiteVisitorCapability
	SendSlackNotification  *SendSlackNotificationCapability
	// TODO above will be deprecated
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

	capabilities := AgentCapabilities{
		AnalyzeWebSession:      NewAnalyzeWebSessionCapability(postgresRepositories, actionService),
		ApplyTag:               NewApplyTagCapability(tagService),
		CreateOrganization:     NewCreateOrganizationCapability(organizationService),
		ICPQualification:       NewICPQualificationCapability(postgresRepositories, aiService),
		IdentifyWebsiteVisitor: NewIdentifyWebsiteVisitorCapability(postgresRepositories, enrichmentService),
		SendSlackNotification:  NewSendSlackNotificationCapability(notificationService),
	}

	executors := make(map[enum.AgentCapabilityType]interfaces.AgentCapabilityUntyped)

	// Register each capability with its corresponding type.
	executors[enum.CapabilityAnalyzeWebSessionIntent] = NewAnalyzeWebSessionCapability(postgresRepositories, actionService)
	executors[enum.CapabilityCreateOrganization] = NewCreateOrganizationCapability(organizationService)
	executors[enum.CapabilityIdentifyWebVisitor] = NewIdentifyWebsiteVisitorCapability(postgresRepositories, enrichmentService)
	executors[enum.CapabilitySendSlackNotification] = NewSendSlackNotificationCapability(notificationService)
	executors[enum.CapabilitySendWebVisitorSlackNotification] = NewSendWebVisitorSlackNotificationCapability(postgresRepositories, notificationService, workspaceService)
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
