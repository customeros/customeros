package agent_capability

import (
	"context"
	"fmt"
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

func (c *IdentifyWebsiteVisitorCapability) ValidateConfig(config NoConfig) error {
	return nil
}

func (c *IdentifyWebsiteVisitorCapability) ValidateInput(data IdentifyWebsiteVisitorInput) error {
	if data.IPAddress == "" {
		return errors.New("IP address cannot be empty")
	}
	if data.SessionID == "" {
		return errors.New("SessionID cannot be empty")
	}
	return nil
}

func (c *IdentifyWebsiteVisitorCapability) GetInput() any {
	return &IdentifyWebsiteVisitorInput{}
}

func (c *IdentifyWebsiteVisitorCapability) GetConfig() any {
	return &NoConfig{}
}

func (c *IdentifyWebsiteVisitorCapability) GetOutput() any {
	return &IdentifyWebsiteVisitorResult{}
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
	_ interfaces.AgentCapabilityExecution[IdentifyWebsiteVisitorInput, IdentifyWebsiteVisitorResult, NoConfig] = (*IdentifyWebsiteVisitorCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                                        = (*IdentifyWebsiteVisitorCapability)(nil)
)

type IdentifyWebsiteVisitorInput struct {
	SessionID string `json:"sessionId"`
	IPAddress string `json:"ipAddress"`
}

type IdentifyWebsiteVisitorResult struct {
	Domain       string `json:"domain"`
	LinkedInSlug string `json:"linkedinSlug"`
}

func (c *IdentifyWebsiteVisitorCapability) Execute(ctx context.Context, data IdentifyWebsiteVisitorInput, config NoConfig) (IdentifyWebsiteVisitorResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)

	results := IdentifyWebsiteVisitorResult{}

	err := c.ValidateInput(data)
	if err != nil {
		tracing.TraceErr(span, err)
		return results, err
	}

	domain, linkedInSlug, err := c.identifyIP(ctx, data.IPAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return results, err
	}
	span.LogKV("domain", domain)

	results.Domain = domain
	results.LinkedInSlug = linkedInSlug

	if domain != "" {
		_, err = c.postgresRepositories.WebSessionRepository.UpdateSessionWithDomain(ctx, data.SessionID, domain)
		if err != nil {
			tracing.TraceErr(span, err)
			return results, err
		}
	}

	return results, nil
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
	if !ok {
		return nil, fmt.Errorf("invalid input type: expected AnalyzeWebSessionInput")
	}

	typedConfig, ok := config.(*NoConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type: expected NoCOnfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
