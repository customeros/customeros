package agent_capability

import (
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
}

func InitCapabilities(
	postgresRepositories *postgres_repository.Repositories,
	actionService interfaces.ActionService,
	aiService interfaces.AIService,
	enrichmentService interfaces.EnrichmentService,
	notificationService interfaces.NotificationService,
	organizationService interfaces.OrganizationService,
	tagService interfaces.TagService,

) *AgentCapabilities {

	capabilities := AgentCapabilities{
		AnalyzeWebSession:      NewAnalyzeWebSessionCapability(postgresRepositories, actionService),
		ApplyTag:               NewApplyTagCapability(tagService),
		CreateOrganization:     NewCreateOrganizationCapability(organizationService),
		ICPQualification:       NewICPQualificationCapability(postgresRepositories, aiService),
		IdentifyWebsiteVisitor: NewIdentifyWebsiteVisitorCapability(postgresRepositories, enrichmentService),
		SendSlackNotification:  NewSendSlackNotificationCapability(notificationService),
	}

	return &capabilities
}
