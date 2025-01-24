package agent_capability

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services"
)

type AgentCapabilities struct {
	AnalyzeWebSession      *AnalyzeWebSessionCapability
	ApplyTag               *ApplyTagCapability
	CreateOrganization     *CreateOrganizationCapability
	ICPQualification       *ICPQualificationCapability
	IdentifyWebsiteVisitor *IdentifyWebsiteVisitorCapability
	SendSlackNotification  *SendSlackNotificationCapability
}

func InitCapabilities(commonServices *service.CommonServices) *AgentCapabilities {

	capabilities := AgentCapabilities{
		AnalyzeWebSession:      NewAnalyzeWebSessionCapability(commonServices.PostgresRepositories, commonServices.ActionService),
		ApplyTag:               NewApplyTagCapability(commonServices.TagService),
		CreateOrganization:     NewCreateOrganizationCapability(commonServices.OrganizationService),
		ICPQualification:       NewICPQualificationCapability(commonServices.PostgresRepositories, commonServices.AIService),
		IdentifyWebsiteVisitor: NewIdentifyWebsiteVisitorCapability(commonServices.PostgresRepositories, commonServices.EnrichmentService),
		SendSlackNotification:  NewSendSlackNotificationCapability(commonServices.NotificationService),
	}

	return &capabilities
}
