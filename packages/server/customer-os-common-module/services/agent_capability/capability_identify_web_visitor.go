package agent_capability

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/opentracing/opentracing-go"
)

type IdentifyWebsiteVisitorInput struct {
	IPAddress string
}

type IdentifyWebsiteVisitorResult struct {
	Domain       string
	LinkedInSlug string
}

func (c *agentCapabilityService) handleIdentifyWebsiteVisitorExecution(ctx context.Context, executionContainer *dto.CapabilityExecutionContainer) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.handleIdentifyWebsiteVisitorExecution")
	defer span.Finish()
	tracing.TagComponentService(span)

	input, ok := executionContainer.InputData.(IdentifyWebsiteVisitorInput)
	if !ok {
		err := fmt.Errorf("expected IdentifyWebsiteVisitorInput, got %T", executionContainer.InputData)
		tracing.TraceErr(span, err)
		return err
	}
	result, err := c.executeIdentifyWebsiteVisitor(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	executionContainer.OutputData = result
	span.LogKV(
		"event", "capability_executed",
		"output_type", fmt.Sprintf("%T", result),
		"domain", result.Domain,
	)
	return nil
}

func (c *agentCapabilityService) executeIdentifyWebsiteVisitor(ctx context.Context, data IdentifyWebsiteVisitorInput) (*IdentifyWebsiteVisitorResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.executeIdentifyWebsiteVisitor")
	defer span.Finish()
	tracing.TagComponentService(span)

	domain, linkedInSlug, err := c.identifyIP(ctx, data.IPAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &IdentifyWebsiteVisitorResult{
		Domain:       domain,
		LinkedInSlug: linkedInSlug,
	}, nil

}

func (c *agentCapabilityService) identifyIP(ctx context.Context, ipAddress string) (domain, linkedinSlug string, err error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentCapabilityService.identifyIP")
	defer span.Finish()

	snitcherData, err := c.enrichmentService.IPIdentity(ctx, ipAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", "", err
	}

	if snitcherData == nil {
		return "", "", nil
	}

	_, primaryDomain := domaincheck.PrimaryDomainCheck(snitcherData.Company.Domain)

	return primaryDomain, snitcherData.Company.Profiles.LinkedIn.Handle, nil
}
