package agent_capability

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type CreateOrganizationCapability struct {
	organizationService interfaces.OrganizationService
}

func NewCreateOrganizationCapability(orgService interfaces.OrganizationService) *CreateOrganizationCapability {
	return &CreateOrganizationCapability{
		organizationService: orgService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[CreateOrganizationInput, CreateOrganizationOutput, NoConfig] = (*CreateOrganizationCapability)(nil)
)

func (c *CreateOrganizationCapability) Type() enum.AgentCapability {
	return enum.CapabilityCreateAndEnrichCompany
}

func (c *CreateOrganizationCapability) GetInput() CreateOrganizationInput {
	return CreateOrganizationInput{}
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

func (c *CreateOrganizationCapability) GetConfig() NoConfig {
	return NoConfig{}
}

func (c *CreateOrganizationCapability) ValidateConfig(config NoConfig) error {
	return nil
}

type CreateOrganizationInput struct {
	Domain string `json:"domain"`
}

type CreateOrganizationOutput struct {
	CapabilityOutput
	OrganizationID string `json:"organizationId"`
}

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
