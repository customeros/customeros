package agent_capability

import (
	"context"
	"errors"

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
var _ interfaces.AgentCapabilityExecution[IdentifyWebsiteVisitorInput, IdentifyWebsiteVisitorResult] = (*IdentifyWebsiteVisitorCapability)(nil)

type IdentifyWebsiteVisitorInput struct {
	SessionID string
	IPAddress string
}

type IdentifyWebsiteVisitorResult struct {
	Domain       string
	LinkedInSlug string
}

func (c *IdentifyWebsiteVisitorCapability) Execute(ctx context.Context, data IdentifyWebsiteVisitorInput) (IdentifyWebsiteVisitorResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IdentifyWebsiteVisitorCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)

	results := IdentifyWebsiteVisitorResult{}

	if data.IPAddress == "" {
		err := errors.New("IP address cannot be empty")
		tracing.TraceErr(span, err)
		return results, err
	}
	if data.SessionID == "" {
		err := errors.New("SessionID cannot be empty")
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
