package agent_capability

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type CreateOrganizationResult struct {
	OrganizationID string
}

func (c *agentCapabilityService) handleOrganizationCreationExecution(ctx context.Context, executionContainer *dto.CapabilityExecutionContainer) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.handleOrganizationCreationExecution")
	defer span.Finish()
	tracing.TagComponentService(span)

	input, ok := executionContainer.InputData.(data_fields.OrganizationFields)
	if !ok {
		err := fmt.Errorf("expected data_fields.OrganizationFields, got %T", executionContainer.InputData)
		tracing.TraceErr(span, err)
		return err
	}

	result, err := c.executeCreateOrganization(ctx, input)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	executionContainer.OutputData = result
	span.LogKV(
		"event", "capability_executed",
		"output_type", fmt.Sprintf("%T", result),
		"id", result.OrganizationID,
	)
	return nil
}

func (c *agentCapabilityService) executeCreateOrganization(ctx context.Context, data data_fields.OrganizationFields) (CreateOrganizationResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.executeCreateOrganization")
	defer span.Finish()
	tracing.TagComponentService(span)

	results := CreateOrganizationResult{}

	orgID, err := c.organizationService.Save(ctx, nil, nil, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return results, err
	}

	results.OrganizationID = orgID

	return results, nil
}
