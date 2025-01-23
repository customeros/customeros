package agent_capability

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
)

type ICPQualificationInput struct {
	ICPDefinition string
	PrimaryDomain string
}

type ICPQualificationResult struct {
}

func (c *agentCapabilityService) handleICPQualificationExecution(ctx context.Context, executionContainer *dto.CapabilityExecutionContainer) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.handleICPQualificationExecution")
	defer span.Finish()
	tracing.TagComponentService(span)

	input, ok := executionContainer.InputData.(ICPQualificationInput)
	if !ok {
		err := fmt.Errorf("expected ICPQualificationInput, got %T", executionContainer.InputData)
		tracing.TraceErr(span, err)
		return err
	}

	_, err := c.executeICPQualification(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (c *agentCapabilityService) executeICPQualification(ctx context.Context, data ICPQualificationInput) (ICPQualificationResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.executeICPQualification")
	defer span.Finish()
	tracing.TagComponentService(span)

	// get ICP definition

	// get all company context
	_, err := c.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, data.PrimaryDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return ICPQualificationResult{}, err
	}

	// build prompt

	// askAI

	// save to timeline

	// disqualify as lead if not a fit

	return ICPQualificationResult{}, nil
}
