package agent_capability

import (
	"context"
	"github.com/opentracing/opentracing-go/log"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type CreateOrganizationCapability struct {
	organizationService interfaces.OrganizationService
	workspaceService    interfaces.WorkspaceService
	events              *events.EventsService
}

func NewCreateOrganizationCapability(
	events *events.EventsService,
	orgService interfaces.OrganizationService,
	workspaceService interfaces.WorkspaceService,
) *CreateOrganizationCapability {
	return &CreateOrganizationCapability{
		events:              events,
		workspaceService:    workspaceService,
		organizationService: orgService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[CreateOrganizationInput, CreateOrganizationOutput, postgres_entity.NoConfig] = (*CreateOrganizationCapability)(nil)
)

func (c *CreateOrganizationCapability) Type() enum.AgentCapability {
	return enum.CapabilityCreateAndEnrichCompany
}

func (c *CreateOrganizationCapability) Name() string {
	return "Create and enrich a companies"
}

func (c *CreateOrganizationCapability) NewInput() CreateOrganizationInput {
	return CreateOrganizationInput{}
}

func (c *CreateOrganizationCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *CreateOrganizationCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
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

func (c *CreateOrganizationCapability) ValidateConfig(config postgres_entity.NoConfig) error {
	return nil
}

type CreateOrganizationInput struct {
	Domain string `json:"domain"`
}

type CreateOrganizationOutput struct {
	OrganizationID    string `json:"organizationId"`
	IsWorkspaceDomain bool   `json:"isWorkspaceDomain"`
}

func (c *CreateOrganizationCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[CreateOrganizationInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, CreateOrganizationOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateOrganizationCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := CreateOrganizationOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	isWorkspaceDomain, err := c.workspaceService.IsWorkspaceDomain(ctx, executionContainer.InputData.Domain)
	if isWorkspaceDomain {
		err = c.events.Publisher.PublishFanoutEvent(ctx, executionContainer.AgentExecutionID, model.AGENT_EXECUTION, dto.WebVisitorNotIdentified{
			AgentExecutionId: executionContainer.AgentExecutionID,
		})
		if err != nil {
			tracing.TraceErr(span, err)
		}

		span.LogFields(log.Bool("result.skip", true))
		result.IsWorkspaceDomain = true
		return enum.CapabilityExecutionCompleted, result, nil
	}

	// Check if the domain is already linked to an organization
	organizationEntity, err := c.organizationService.GetOrganizationByDomain(ctx, executionContainer.InputData.Domain, true)

	// organization found
	if organizationEntity != nil {
		// if organization is hidden, return early
		if organizationEntity.Hide {
			return enum.CapabilityExecutionError, result, errors.New("Identified organization is archived")
		}
		result.OrganizationID = organizationEntity.ID
		tracing.LogObjectAsJson(span, "result", result)
		return enum.CapabilityExecutionCompleted, result, nil
	}

	orgID, err := c.organizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
		Domains: []string{executionContainer.InputData.Domain},
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	err = c.publishCompanyIdentifiedEvent(ctx, executionContainer.AgentExecutionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	result.OrganizationID = orgID
	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}

func (c *CreateOrganizationCapability) publishCompanyIdentifiedEvent(ctx context.Context, agentExecutionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateOrganizationCapability.publishCompanyIdentifiedEvent")
	defer span.Finish()
	tracing.TagComponentService(span)

	return c.events.Publisher.PublishFanoutEvent(ctx, agentExecutionID, model.AGENT_EXECUTION, dto.CompanyIdentified{
		AgentExecutionId: agentExecutionID,
	})
}
