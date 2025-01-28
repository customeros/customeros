package agent_capability

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"

	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type IdentifyWebsiteVisitorCapability struct {
	postgresRepositories *postgres_repository.Repositories
	enrichmentService    interfaces.EnrichmentService
}

type IdentifyWebsiteVisitorInput struct {
	SessionID string `json:"sessionId"`
	IPAddress string `json:"ipAddress"`
	VisitorID string `json:"visitorId"`
	Hostname  string `json:"hostname"`
}

type IdentifyWebsiteVisitorOutput struct {
	CapabilityOutput
	Domain       string `json:"domain"`
	LinkedInSlug string `json:"linkedinSlug"`
}

type IdentifyWebsiteVisitorConfig struct {
	Websites WebsitesConfig `json:"websites"`
}

type WebsitesConfig struct {
	Value []string `json:"value"`
	Error string   `json:"error"`
}

func (c *IdentifyWebsiteVisitorCapability) ValidateConfig(config IdentifyWebsiteVisitorConfig) error {
	return nil
}

func (c *IdentifyWebsiteVisitorCapability) ValidateInput(data IdentifyWebsiteVisitorInput) error {
	if data.IPAddress == "" {
		return errors.New("IP address cannot be empty")
	}
	if data.SessionID == "" {
		return errors.New("SessionID cannot be empty")
	}
	if data.VisitorID == "" {
		return errors.New("VisitorID cannot be empty")
	}
	return nil
}

func (c *IdentifyWebsiteVisitorCapability) GetInput() any {
	return &IdentifyWebsiteVisitorInput{}
}

func (c *IdentifyWebsiteVisitorCapability) GetConfig() any {
	return &IdentifyWebsiteVisitorConfig{}
}

func (c *IdentifyWebsiteVisitorCapability) GetOutput() any {
	return &IdentifyWebsiteVisitorOutput{}
}

func NewIdentifyWebsiteVisitorCapability(
	postgresRepositories *postgres_repository.Repositories,
	enrichmentService interfaces.EnrichmentService,
) *IdentifyWebsiteVisitorCapability {
	return &IdentifyWebsiteVisitorCapability{
		postgresRepositories: postgresRepositories,
		enrichmentService:    enrichmentService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapabilityExecution[IdentifyWebsiteVisitorInput, IdentifyWebsiteVisitorOutput, IdentifyWebsiteVisitorConfig] = (*IdentifyWebsiteVisitorCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                                                            = (*IdentifyWebsiteVisitorCapability)(nil)
)

func (c *IdentifyWebsiteVisitorCapability) Execute(ctx context.Context, data IdentifyWebsiteVisitorInput, config IdentifyWebsiteVisitorConfig) (IdentifyWebsiteVisitorOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", data)
	tracing.LogObjectAsJson(span, "config", config)

	result := IdentifyWebsiteVisitorOutput{}

	if err := c.ValidateInput(data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return result, err
	}
	if err := c.ValidateConfig(config); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return result, err
	}

	// check if website input hostname configured by current agent
	err := c.acceptHostname(ctx, data.Hostname, config.Websites.Value)
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}

	result.ExecutionValidated = true

	domain, linkedInSlug, err := c.identifyIP(ctx, data.IPAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}
	span.LogKV("domain", domain)

	result.Domain = domain
	result.LinkedInSlug = linkedInSlug

	if domain != "" {
		_, err = c.postgresRepositories.WebSessionRepository.UpdateSessionWithDomain(ctx, data.SessionID, domain)
		if err != nil {
			tracing.TraceErr(span, err)
			return result, err
		}
	}

	result.Completed = true
	tracing.LogObjectAsJson(span, "result", result)
	return result, nil
}

func (c *IdentifyWebsiteVisitorCapability) acceptHostname(ctx context.Context, hostname string, websites []string) error {
	span, ctx := tracing.StartTracerSpan(ctx, "IdentifyWebsiteVisitorCapability.acceptHostname")
	defer span.Finish()

	accepted := false

	for _, website := range websites {
		if utils.CleanUrlBasePath(website) == utils.CleanUrlBasePath(hostname) {
			accepted = true
		}
	}

	if !accepted {
		return coserrors.ErrCapabilityHostnameNotConfigured
	}

	return nil
}

func (c *IdentifyWebsiteVisitorCapability) identifyIP(ctx context.Context, ipAddress string) (domain, linkedinSlug string, err error) {
	span, ctx := tracing.StartTracerSpan(ctx, "IdentifyWebsiteVisitorCapability.identifyIP")
	defer span.Finish()

	snitcherData, err := c.enrichmentService.IPIdentity(ctx, ipAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", "", err
	}

	if snitcherData == nil || snitcherData.Company == nil || snitcherData.Company.Domain == "" {
		return "", "", nil
	}

	_, primaryDomain := domaincheck.PrimaryDomainCheck(snitcherData.Company.Domain)

	if snitcherData.Company.Profiles != nil && snitcherData.Company.Profiles.LinkedIn != nil {
		linkedinSlug = snitcherData.Company.Profiles.LinkedIn.Handle
	}

	return primaryDomain, snitcherData.Company.Profiles.LinkedIn.Handle, nil
}

// ExecuteUntyped implements the AgentCapabilityUntyped interface.
// It casts the generic input and config to the specific types and delegates to the typed Execute method.
func (c *IdentifyWebsiteVisitorCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
	typedInput, ok := input.(*IdentifyWebsiteVisitorInput)
	if !ok || typedInput == nil {
		return nil, fmt.Errorf("invalid input type: expected AnalyzeWebSessionInput")
	}

	typedConfig, ok := config.(*IdentifyWebsiteVisitorConfig)
	if !ok || typedConfig == nil {
		return nil, fmt.Errorf("invalid config type: expected IdentifyWebsiteVisitorConfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
