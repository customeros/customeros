package agent_capability

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type CreateOrganizationCapability struct {
	organizationService interfaces.OrganizationService
}

func NewCreateOrganizationCapability(orgService interfaces.OrganizationService) *CreateOrganizationCapability {
	return &CreateOrganizationCapability{
		organizationService: orgService,
	}
}

// Compile-time interface check
var _ interfaces.AgentCapabilityExecution[CreateOrganizationInput, CreateOrganizationResult] = (*CreateOrganizationCapability)(nil)

type CreateOrganizationInput struct {
	Domain string
}

type CreateOrganizationResult struct {
	OrganizationID string
}

func (c *CreateOrganizationCapability) Execute(ctx context.Context, data CreateOrganizationInput) (CreateOrganizationResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateOrganizationCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)

	results := CreateOrganizationResult{}

	orgID, err := c.organizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
		Domains: []string{data.Domain},
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return results, err
	}

	results.OrganizationID = orgID

	return results, nil
}
