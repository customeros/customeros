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
	contactService interfaces.ContactService,
	domainService interfaces.DomainService,
	enrichmentService interfaces.EnrichmentService,
	invoiceService interfaces.InvoiceService,
	markdownService interfaces.MarkdownEventService,
	notificationService interfaces.NotificationService,
	organizationService interfaces.OrganizationService,
	tagService interfaces.TagService,
	workspaceService interfaces.WorkspaceService,
) *AgentCapabilities {
	var capabilities []interfaces.AgentCapabilityUntyped
	capabilities = append(capabilities, NewAddMeetingNotesToCompanyCapability(organizationService, markdownService))
	capabilities = append(capabilities, NewAnalyzeWebSessionCapability(events, postgresRepositories, actionService))
	capabilities = append(capabilities, NewApplyTagCapability(tagService))
	capabilities = append(capabilities, NewCreateOrganizationCapability(organizationService))
	capabilities = append(capabilities, NewCreateContactCapability(contactService))
	capabilities = append(capabilities, NewCreateMarkdownTimelineEventCapability(markdownService))
	capabilities = append(capabilities, NewEvaluateICPFitCapability(aiService))
	capabilities = append(capabilities, NewExtractMeetingHighlightsCapability(aiService))
	capabilities = append(capabilities, NewGatherCompanyIntelligenceCapability(postgresRepositories, organizationService))
	capabilities = append(capabilities, NewGenerateInvoiceCapability(postgresRepositories, invoiceService))
	capabilities = append(capabilities, NewIdentifyMeetingParticipantsCapability(workspaceService))
	capabilities = append(capabilities, NewIdentifyWebsiteVisitorCapability(events, postgresRepositories, enrichmentService, domainService))
	capabilities = append(capabilities, NewSendSlackNotificationCapability(notificationService))
	capabilities = append(capabilities, NewSendWebVisitorSlackNotificationCapability(postgresRepositories, notificationService, workspaceService))
	capabilities = append(capabilities, NewUpdateCompanyStatusCapability(organizationService, events))

	agentCapabilities := AgentCapabilities{
		executors: make(map[enum.AgentCapability]interfaces.AgentCapabilityUntyped),
	}

	for _, capability := range capabilities {
		agentCapabilities.executors[capability.Type()] = capability
	}

	return &agentCapabilities
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

	capability, ok := executor.(interfaces.AgentCapability[I, O, C])
	if !ok {
		return nil, false
	}
	return capability, true
}
