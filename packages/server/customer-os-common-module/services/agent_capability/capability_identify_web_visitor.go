package agent_capability

import (
	"context"
	"time"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

const (
	IPHistoryTTL = 30 * 24 * time.Hour
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

func (c *IdentifyWebsiteVisitorCapability) DefaultActive() bool {
	return true
}

func (c *IdentifyWebsiteVisitorCapability) ValidateInput(data IdentifyWebsiteVisitorInput) error {
	if data.WebSessionID == "" {
		return errors.New("WebSessionID cannot be empty")
	}
	return nil
}

func (c *IdentifyWebsiteVisitorCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

type IdentifyWebsiteVisitorInput struct {
	WebSessionID string `json:"webSessionId"`
}

type IdentifyWebsiteVisitorOutput struct {
	Domain       string `json:"domain"`
	EmailAddress string `json:"emailAddress"`
}

func (c *IdentifyWebsiteVisitorCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[IdentifyWebsiteVisitorInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, IdentifyWebsiteVisitorOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.Execute")
	defer span.Finish()
	tracing.SetDefaultAgentCapabilitySpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "executionContainer", executionContainer)

	result := IdentifyWebsiteVisitorOutput{}

	// Validate input and config
	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	// Get web session
	webSession, err := c.getWebSession(ctx, executionContainer.InputData.WebSessionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	// Check if identity already set
	identityStatus, err := c.processExistingIdentity(ctx, webSession, executionContainer.AgentExecutionID, &result)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	// Return if identity is already set or if we need to stop execution
	if identityStatus != nil {
		return *identityStatus, result, nil
	}

	// Check if we have identifiable information from other sessions with the same IP
	identityFromIPStatus, err := c.tryIdentifyFromIPHistory(ctx, webSession, executionContainer.AgentExecutionID, &result)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	// Return if we found identity from IP history
	if identityFromIPStatus != nil {
		tracing.LogObjectAsJson(span, "result", result)
		return *identityFromIPStatus, result, nil
	}

	// Attempt to identify using third-party enrichment services
	status, output, err := c.tryIdentifyFromEnrichment(ctx, webSession, executionContainer.AgentExecutionID, &result)
	tracing.LogObjectAsJson(span, "result", output)
	return status, output, err
}

// getWebSession retrieves the web session with the given ID
func (c *IdentifyWebsiteVisitorCapability) getWebSession(ctx context.Context, sessionID string) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.getWebSession")
	defer span.Finish()
	tracing.SetDefaultAgentCapabilitySpanTags(ctx, span)

	webSession, err := c.postgresRepositories.WebSessionRepository.FindSession(ctx, postgres_entity.WebSession{
		ID:     sessionID,
		Tenant: common.GetTenantFromContext(ctx),
	}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find web session")
	}
	if webSession == nil {
		return nil, errors.New("web session not found")
	}

	return webSession, nil
}

// processExistingIdentity checks if identity is already set on the web session
// Returns nil if identity is not set and processing should continue
// Returns a status if processing should stop
func (c *IdentifyWebsiteVisitorCapability) processExistingIdentity(
	ctx context.Context,
	websession *postgres_entity.WebSession,
	agentExecutionID string,
	result *IdentifyWebsiteVisitorOutput,
) (*enum.CapabilityExecutionStatus, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.processExistingIdentity")
	defer span.Finish()
	tracing.SetDefaultAgentCapabilitySpanTags(ctx, span)

	// Check if identity already set
	identitySet, primaryDomain, err := c.isIdentitySetOnWebSession(ctx, *websession)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check if identity is set")
	}

	// If identity not set, continue with processing
	if !identitySet {
		return nil, nil
	}

	// Check if it's a workspace domain
	isWorkspaceDomain, err := c.workspaceService.IsWorkspaceDomain(ctx, primaryDomain)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check if domain is workspace domain")
	}

	// If it's a workspace domain, publish event and stop execution
	if isWorkspaceDomain {
		err = c.publishWebVisitorNotIdentifiedEvent(ctx, agentExecutionID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to publish web visitor not identified event")
		}

		status := enum.CapabilityExecutionStop
		return &status, nil
	}

	// Identity found and it's not a workspace domain
	result.Domain = primaryDomain
	if websession.Email != nil {
		result.EmailAddress = *websession.Email
	}

	err = c.publishWebVisitorIdentifiedEvent(ctx, agentExecutionID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to publish web visitor identified event")
	}

	status := enum.CapabilityExecutionCompleted
	return &status, nil
}

// tryIdentifyFromIPHistory attempts to identify a visitor using historical sessions with the same IP
// Returns nil if identification was not possible and processing should continue
func (c *IdentifyWebsiteVisitorCapability) tryIdentifyFromIPHistory(
	ctx context.Context,
	websession *postgres_entity.WebSession,
	agentExecutionID string,
	result *IdentifyWebsiteVisitorOutput,
) (*enum.CapabilityExecutionStatus, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.tryIdentifyFromIPHistory")
	defer span.Finish()
	tracing.SetDefaultAgentCapabilitySpanTags(ctx, span)

	// Find session with same IP that has domain information
	identifiedSession, err := c.postgresRepositories.WebSessionRepository.FindLatestSessionWithDomainByIP(ctx, websession.IP)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find session with domain by IP")
	}

	// If no identified session found, continue with processing
	if identifiedSession == nil || identifiedSession.Domain == nil || *identifiedSession.Domain == "" {
		return nil, nil
	}
	// If identified session is older that TTL, continue with processing
	if time.Since(identifiedSession.LastActivity) > IPHistoryTTL {
		return nil, nil
	}

	// Identified from IP history
	result.Domain = utils.IfNotNilString(identifiedSession.Domain)
	if utils.IfNotNilString(identifiedSession.Email) != "" {
		result.EmailAddress = utils.IfNotNilString(identifiedSession.Email)
	}

	// update current session with identified domain
	err = c.postgresRepositories.WebSessionRepository.SetVisitorIdentity(ctx, websession.ID, identifiedSession.Domain, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set visitor identity")
	}

	err = c.publishWebVisitorIdentifiedEvent(ctx, agentExecutionID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to publish web visitor identified event")
	}

	status := enum.CapabilityExecutionCompleted
	return &status, nil
}

// tryIdentifyFromEnrichment attempts to identify a visitor using third-party enrichment services
func (c *IdentifyWebsiteVisitorCapability) tryIdentifyFromEnrichment(
	ctx context.Context,
	websession *postgres_entity.WebSession,
	agentExecutionID string,
	result *IdentifyWebsiteVisitorOutput,
) (enum.CapabilityExecutionStatus, IdentifyWebsiteVisitorOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.tryIdentifyFromEnrichment")
	defer span.Finish()
	tracing.SetDefaultAgentCapabilitySpanTags(ctx, span)

	// Try to identify IP via 3rd parties
	domain, err := c.identifyIP(ctx, websession.IP)
	if err != nil {
		return enum.CapabilityExecutionError, *result, errors.Wrap(err, "failed to identify IP")
	}

	result.Domain = domain

	// If domain not found, publish event and complete
	if domain == "" {
		err = c.publishWebVisitorNotIdentifiedEvent(ctx, agentExecutionID)
		if err != nil {
			return enum.CapabilityExecutionError, *result, errors.Wrap(err, "failed to publish web visitor not identified event")
		}
		return enum.CapabilityExecutionCompleted, *result, nil
	}

	// Update the web session with identified domain
	err = c.postgresRepositories.WebSessionRepository.SetVisitorIdentity(ctx, websession.ID, &domain, nil, nil)
	if err != nil {
		return enum.CapabilityExecutionError, *result, errors.Wrap(err, "failed to set visitor identity")
	}

	// Check if it's a workspace domain
	isWorkspaceDomain, err := c.workspaceService.IsWorkspaceDomain(ctx, domain)
	if err != nil {
		return enum.CapabilityExecutionError, *result, errors.Wrap(err, "failed to check if domain is workspace domain")
	}

	if isWorkspaceDomain {
		err = c.publishWebVisitorNotIdentifiedEvent(ctx, agentExecutionID)
		if err != nil {
			return enum.CapabilityExecutionError, *result, errors.Wrap(err, "failed to publish web visitor not identified event")
		}
		return enum.CapabilityExecutionStop, *result, nil
	}

	// Publish success event
	err = c.publishWebVisitorIdentifiedEvent(ctx, agentExecutionID)
	if err != nil {
		return enum.CapabilityExecutionError, *result, errors.Wrap(err, "failed to publish web visitor identified event")
	}

	return enum.CapabilityExecutionCompleted, *result, nil
}

// isIdentitySetOnWebSession checks if identity information is already set on the web session
func (c *IdentifyWebsiteVisitorCapability) isIdentitySetOnWebSession(ctx context.Context, websession postgres_entity.WebSession) (bool, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.isIdentitySetOnWebSession")
	defer span.Finish()
	tracing.SetDefaultAgentCapabilitySpanTags(ctx, span)

	if websession.Domain == nil || *websession.Domain == "" {
		return false, "", nil
	}

	_, _, primaryDomain := c.domainService.CheckDomainWithMailsherpa(ctx, *websession.Domain)

	// Update session with primary domain if it differs from current domain
	if primaryDomain != *websession.Domain {
		err := c.postgresRepositories.WebSessionRepository.SetVisitorIdentity(ctx, websession.ID, &primaryDomain, nil, nil)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	return true, primaryDomain, nil
}

// publishWebVisitorIdentifiedEvent publishes an event indicating a visitor was identified
func (c *IdentifyWebsiteVisitorCapability) publishWebVisitorIdentifiedEvent(ctx context.Context, agentExecutionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.publishWebVisitorIdentifiedEvent")
	defer span.Finish()
	tracing.SetDefaultAgentCapabilitySpanTags(ctx, span)

	return c.events.Publisher.PublishFanoutEvent(ctx, agentExecutionID, model.AGENT_EXECUTION, dto.WebVisitorIdentified{
		AgentExecutionId: agentExecutionID,
	})
}

// publishWebVisitorNotIdentifiedEvent publishes an event indicating a visitor was not identified
func (c *IdentifyWebsiteVisitorCapability) publishWebVisitorNotIdentifiedEvent(ctx context.Context, agentExecutionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.publishWebVisitorNotIdentifiedEvent")
	defer span.Finish()
	tracing.SetDefaultAgentCapabilitySpanTags(ctx, span)

	return c.events.Publisher.PublishFanoutEvent(ctx, agentExecutionID, model.AGENT_EXECUTION, dto.WebVisitorNotIdentified{
		AgentExecutionId: agentExecutionID,
	})
}

// identifyIP attempts to identify a company domain from an IP address using enrichment services
func (c *IdentifyWebsiteVisitorCapability) identifyIP(ctx context.Context, ipAddress string) (primaryDomain string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.identifyIP")
	defer span.Finish()
	span.LogKV("ipAddress", ipAddress)

	snitcherData, err := c.enrichmentService.IPIdentity(ctx, ipAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "failed to get IP identity")
	}

	if snitcherData == nil || snitcherData.Company == nil || snitcherData.Company.Domain == "" {
		return "", nil
	}

	_, _, primaryDomain = c.domainService.CheckDomainWithMailsherpa(ctx, snitcherData.Company.Domain)

	return primaryDomain, nil
}
