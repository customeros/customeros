package agent_capability

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type IdentifyWebsiteVisitorCapability struct {
	events               *events.EventsService
	postgresRepositories *postgres_repository.Repositories
	enrichmentService    interfaces.EnrichmentService
	domainService        interfaces.DomainService
	workspaceService     interfaces.WorkspaceService
}

func NewIdentifyWebsiteVisitorCapability(
	events *events.EventsService,
	postgresRepositories *postgres_repository.Repositories,
	enrichmentService interfaces.EnrichmentService,
	domainService interfaces.DomainService,
	workspaceService interfaces.WorkspaceService,
) *IdentifyWebsiteVisitorCapability {
	return &IdentifyWebsiteVisitorCapability{
		events:               events,
		postgresRepositories: postgresRepositories,
		enrichmentService:    enrichmentService,
		domainService:        domainService,
		workspaceService:     workspaceService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[IdentifyWebsiteVisitorInput, IdentifyWebsiteVisitorOutput, postgres_entity.NoConfig] = (*IdentifyWebsiteVisitorCapability)(nil)
)

func (c *IdentifyWebsiteVisitorCapability) Type() enum.AgentCapability {
	return enum.CapabilityIdentifyWebVisitor
}

func (c *IdentifyWebsiteVisitorCapability) Name() string {
	return "Identify website visitors"
}

func (c *IdentifyWebsiteVisitorCapability) NewInput() IdentifyWebsiteVisitorInput {
	return IdentifyWebsiteVisitorInput{}
}

func (c *IdentifyWebsiteVisitorCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *IdentifyWebsiteVisitorCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *IdentifyWebsiteVisitorCapability) ValidateInput(data IdentifyWebsiteVisitorInput) error {
	if data.IPAddress == "" {
		return errors.New("IP address cannot be empty")
	}
	if data.WebSessionID == "" {
		return errors.New("WebSessionID cannot be empty")
	}
	if data.VisitorID == "" {
		return errors.New("VisitorID cannot be empty")
	}
	return nil
}

func (c *IdentifyWebsiteVisitorCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

type IdentifyWebsiteVisitorInput struct {
	WebSessionID string `json:"webSessionId"`
	IPAddress    string `json:"ipAddress"`
	VisitorID    string `json:"visitorId"`
	Hostname     string `json:"hostname"`
}

type IdentifyWebsiteVisitorOutput struct {
	Domain       string `json:"domain"`
	LinkedInSlug string `json:"linkedinSlug"`
}

func (c *IdentifyWebsiteVisitorCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[IdentifyWebsiteVisitorInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, IdentifyWebsiteVisitorOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "executionContainer", executionContainer)

	result := IdentifyWebsiteVisitorOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	domain, linkedInSlug, err := c.identifyIP(ctx, executionContainer.InputData.IPAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		tracing.LogObjectAsJson(span, "result", result)
		return enum.CapabilityExecutionError, result, err
	}
	span.LogKV("domain", domain)

	result.Domain = domain
	result.LinkedInSlug = linkedInSlug
	tracing.LogObjectAsJson(span, "result", result)

	if domain != "" {
		_, err = c.postgresRepositories.WebSessionRepository.UpdateSessionWithDomain(ctx, executionContainer.InputData.WebSessionID, domain)
		if err != nil {
			tracing.TraceErr(span, err)
			tracing.LogObjectAsJson(span, "result", result)
			return enum.CapabilityExecutionError, result, err
		}

		isWorkspaceDomain, err := c.workspaceService.IsWorkspaceDomain(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, err)
			return enum.CapabilityExecutionError, result, err
		}
		if isWorkspaceDomain {
			err = c.publishWebVisitorNotIdentifiedEvent(ctx, executionContainer.AgentExecutionID)
			if err != nil {
				tracing.TraceErr(span, err)
				return enum.CapabilityExecutionError, result, err
			}
			return enum.CapabilityExecutionCompleted, result, nil
		}

		err = c.publishWebVisitorIdentifiedEvent(ctx, executionContainer.AgentExecutionID)
		if err != nil {
			tracing.TraceErr(span, err)
			return enum.CapabilityExecutionError, result, err
		}

		return enum.CapabilityExecutionCompleted, result, nil
	}

	err = c.publishWebVisitorNotIdentifiedEvent(ctx, executionContainer.AgentExecutionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	return enum.CapabilityExecutionCompleted, result, nil
}

func (c *IdentifyWebsiteVisitorCapability) publishWebVisitorIdentifiedEvent(ctx context.Context, agentExecutionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.publishWebVisitorIdentifiedEvent")
	defer span.Finish()
	tracing.TagComponentService(span)

	return c.events.Publisher.PublishFanoutEvent(ctx, agentExecutionID, model.AGENT_EXECUTION, dto.WebVisitorIdentified{
		AgentExecutionId: agentExecutionID,
	})
}

func (c *IdentifyWebsiteVisitorCapability) publishWebVisitorNotIdentifiedEvent(ctx context.Context, agentExecutionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.publishWebVisitorNotIdentifiedEvent")
	defer span.Finish()
	tracing.TagComponentService(span)

	return c.events.Publisher.PublishFanoutEvent(ctx, agentExecutionID, model.AGENT_EXECUTION, dto.WebVisitorNotIdentified{
		AgentExecutionId: agentExecutionID,
	})
}

func (c *IdentifyWebsiteVisitorCapability) acceptHostname(ctx context.Context, hostname string, websites []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.acceptHostname")
	defer span.Finish()
	tracing.TagComponentService(span)

	accepted := false

	for _, website := range websites {
		if utils.StripUrlToBasePath(website) == utils.StripUrlToBasePath(hostname) {
			accepted = true
		}
	}

	if !accepted {
		return coserrors.ErrCapabilityHostnameNotConfigured
	}

	return nil
}

func (c *IdentifyWebsiteVisitorCapability) identifyIP(ctx context.Context, ipAddress string) (primaryDomain, linkedinSlug string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.identifyIP")
	defer span.Finish()
	span.LogKV("ipAddress", ipAddress)

	snitcherData, err := c.enrichmentService.IPIdentity(ctx, ipAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", "", err
	}

	if snitcherData == nil || snitcherData.Company == nil || snitcherData.Company.Domain == "" {
		return "", "", nil
	}

	_, _, primaryDomain = c.domainService.CheckDomainWithMailsherpa(ctx, snitcherData.Company.Domain)

	if snitcherData.Company.Profiles != nil && snitcherData.Company.Profiles.LinkedIn != nil {
		linkedinSlug = snitcherData.Company.Profiles.LinkedIn.Handle
	}

	return primaryDomain, snitcherData.Company.Profiles.LinkedIn.Handle, nil
}
