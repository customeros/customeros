package agent_capability

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type CreateOrganizationCapability struct {
	organizationService interfaces.OrganizationService
	domainService       interfaces.DomainService
}

type CreateOrganizationInput struct {
	Domain string `json:"domain"`
}

type CreateOrganizationOutput struct {
	CapabilityOutput
	OrganizationID string `json:"organizationId"`
}

func NewCreateOrganizationCapability(orgService interfaces.OrganizationService, domainService interfaces.DomainService) *CreateOrganizationCapability {
	return &CreateOrganizationCapability{
		organizationService: orgService,
		domainService:       domainService,
	}
}

func (c *CreateOrganizationCapability) ValidateConfig(config NoConfig) error {
	return nil
}

func (c *CreateOrganizationCapability) ValidateInput(input CreateOrganizationInput) error {
	if input.Domain == "" {
		return coserrors.ErrCapabilityDomainMissing
	}
	if !utils.IsValidDomain(input.Domain) {
		return errors.New("Invalid domain format")
	}
	return nil
}

func (c *CreateOrganizationCapability) GetInput() any {
	return &CreateOrganizationInput{}
}

func (c *CreateOrganizationCapability) GetConfig() any {
	return &NoConfig{}
}

func (c *CreateOrganizationCapability) GetOutput() any {
	return &CreateOrganizationOutput{}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapabilityExecution[CreateOrganizationInput, CreateOrganizationOutput, NoConfig] = (*CreateOrganizationCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                                = (*CreateOrganizationCapability)(nil)
)

func (c *CreateOrganizationCapability) Execute(ctx context.Context, data CreateOrganizationInput, config NoConfig) (CreateOrganizationOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateOrganizationCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", data)
	tracing.LogObjectAsJson(span, "config", config)

	result := CreateOrganizationOutput{}

	if err := c.ValidateInput(data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return result, err
	}
	if err := c.ValidateConfig(config); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return result, err
	}

	result.ExecutionValidated = true

	// Check if the domain is already linked to an organization
	organizationEntity, err := c.organizationService.GetOrganizationByDomain(ctx, data.Domain, true)

	// organization found
	if organizationEntity != nil {
		// if organization is hidden, return early
		if organizationEntity.Hide {
			return result, errors.New("Identified organization is archived")
		}
		result.OrganizationID = organizationEntity.ID
		result.Completed = true
		tracing.LogObjectAsJson(span, "result", result)
		return result, nil
	}

	orgID, err := c.organizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
		Domains: []string{data.Domain},
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}

	result.OrganizationID = orgID
	result.Completed = true
	tracing.LogObjectAsJson(span, "result", result)
	return result, nil
}

// ExecuteUntyped implements the AgentCapabilityUntyped interface.
// It casts the generic input and config to the specific types and delegates to the typed Execute method.
func (c *CreateOrganizationCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
	typedInput, ok := input.(*CreateOrganizationInput)
	if !ok || typedInput == nil {
		return nil, fmt.Errorf("invalid input type for CreateOrganizationCapability: expected CreateOrganizationInput")
	}

	typedConfig, ok := config.(*NoConfig)
	if !ok || typedConfig == nil {
		return nil, fmt.Errorf("invalid config type for CreateOrganizationCapability: expected CreateOrganizationConfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
