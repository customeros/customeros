package agent_capability

import (
	"fmt"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
)

type AgentCapabilities struct {
	executors map[enum.AgentCapability]interfaces.AgentCapabilityUntyped
}

func InitCapabilities(
	events *events.EventsService,
	postgresRepositories *postgres_repository.Repositories,
	actionService interfaces.ActionService,
	aiService interfaces.AIService,
	enrichmentService interfaces.EnrichmentService,
	notificationService interfaces.NotificationService,
	organizationService interfaces.OrganizationService,
	tagService interfaces.TagService,
	workspaceService interfaces.WorkspaceService,
	markdownService interfaces.MarkdownEventService,
	domainService interfaces.DomainService,
	invoiceService interfaces.InvoiceService,
) *AgentCapabilities {
	capabilities := &AgentCapabilities{
		executors: make(map[enum.AgentCapability]interfaces.AgentCapabilityUntyped),
	}

	// Register each capability with its corresponding type.
	capabilities.executors[enum.CapabilityAnalyzeWebSessionIntent] = NewAnalyzeWebSessionCapability(events, postgresRepositories, actionService)
	capabilities.executors[enum.CapabilityApplyTag] = NewApplyTagCapability(tagService)
	capabilities.executors[enum.CapabilityCreateAndEnrichCompany] = NewCreateOrganizationCapability(organizationService)
	capabilities.executors[enum.CapabilityCreateMarkdownTimelineEvent] = NewCreateMarkdownTimelineEventCapability(markdownService)
	capabilities.executors[enum.CapabilityEvaluateCompanyICPFit] = NewEvaluateICPFitCapability(aiService, events)
	capabilities.executors[enum.CapabilityGatherCompanyIntelligence] = NewGatherCompanyIntelligenceCapability(postgresRepositories, organizationService)
	capabilities.executors[enum.CapabilityGenerateInvoice] = NewGenerateInvoiceCapability(postgresRepositories, invoiceService)
	capabilities.executors[enum.CapabilityIdentifyWebVisitor] = NewIdentifyWebsiteVisitorCapability(events, postgresRepositories, enrichmentService, domainService)
	capabilities.executors[enum.CapabilitySendSlackNotification] = NewSendSlackNotificationCapability(notificationService)
	capabilities.executors[enum.CapabilitySendWebVisitorSlackNotification] = NewSendWebVisitorSlackNotificationCapability(postgresRepositories, notificationService, workspaceService)
	capabilities.executors[enum.CapabilityUpdateCompanyStatus] = NewUpdateCompanyStatusCapability(organizationService)
	// Continue registering other capabilities here...

	return capabilities
}

// GetExecutor retrieves the untyped executor based on the capability type.
// Returns an error if the capability type is unsupported.
func (c *AgentCapabilities) GetExecutor(capType enum.AgentCapability) (interfaces.AgentCapabilityUntyped, error) {
	executor, exists := c.executors[capType]
	if !exists {
		return nil, fmt.Errorf("unsupported capability type: %s", capType)
	}
	return executor, nil
}

func (c *AgentCapabilities) GetExecutors() map[enum.AgentCapability]interfaces.AgentCapabilityUntyped {
	return c.executors
}

func GetTypedExecutor[I, O, C any](
	executors map[enum.AgentCapability]interfaces.AgentCapabilityUntyped,
	capType enum.AgentCapability,
) (interfaces.AgentCapability[I, O, C], bool) {
	key := capType
	executor, exists := executors[key]
	if !exists {
		return nil, false
	}

	cap, ok := executor.(interfaces.AgentCapability[I, O, C])
	if !ok {
		return nil, false
	}
	return cap, true
}
