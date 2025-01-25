package agent_capability

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type CreateOrganizationCapability struct {
	organizationService interfaces.OrganizationService
}

func (c *CreateOrganizationCapability) GetInput() any {
	return &CreateOrganizationInput{}
}

func (c *CreateOrganizationCapability) GetConfig() any {
	return &NoConfig{}
}

func (c *CreateOrganizationCapability) GetOutput() any {
	return &CreateOrganizationResult{}
}

// CapType returns the capability's corresponding AgentCapabilityType.
func (c *CreateOrganizationCapability) CapType() enum.AgentCapabilityType {
	return enum.CapabilityCreateOrganization
}

func NewCreateOrganizationCapability(orgService interfaces.OrganizationService) *CreateOrganizationCapability {
	return &CreateOrganizationCapability{
		organizationService: orgService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapabilityExecution[CreateOrganizationInput, CreateOrganizationResult, NoConfig] = (*CreateOrganizationCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                                = (*CreateOrganizationCapability)(nil)
)

type CreateOrganizationInput struct {
	Domain string `json:"domain"`
}

type CreateOrganizationResult struct {
	OrganizationID string `json:"organizationId"`
}

func (c *CreateOrganizationCapability) Execute(ctx context.Context, data CreateOrganizationInput, config NoConfig) (CreateOrganizationResult, error) {
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

// ExecuteUntyped implements the AgentCapabilityUntyped interface.
// It casts the generic input and config to the specific types and delegates to the typed Execute method.
func (c *CreateOrganizationCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
	typedInput, ok := input.(*CreateOrganizationInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for CreateOrganizationCapability: expected CreateOrganizationInput")
	}

	typedConfig, ok := config.(*NoConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type for CreateOrganizationCapability: expected CreateOrganizationConfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
