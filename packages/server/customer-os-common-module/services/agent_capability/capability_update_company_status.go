package agent_capability

import (
	"context"

	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type UpdateCompanyStatusCapability struct {
	organizationService interfaces.OrganizationService
	events              *events.EventsService
}

type UpdateCompanyStatusInput struct {
	OrganizationID           string      `json:"organizationId"`
	OrganizationRelationship string      `json:"organizationRelationship"`
	OrganizationStage        string      `json:"organizationStage"`
	IcpFit                   enum.IcpFit `json:"icpFit"`
	IcpFitRationale          []string    `json:"icpFitRationale"`
}

func NewUpdateCompanyStatusCapability(orgSrv interfaces.OrganizationService, events *events.EventsService) *UpdateCompanyStatusCapability {
	return &UpdateCompanyStatusCapability{
		organizationService: orgSrv,
		events:              events,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[UpdateCompanyStatusInput, NoOutput, postgres_entity.NoConfig] = (*UpdateCompanyStatusCapability)(nil)
)

func (c *UpdateCompanyStatusCapability) Type() enum.AgentCapability {
	return enum.CapabilityUpdateCompanyStatus
}

func (c *UpdateCompanyStatusCapability) Name() string {
	return "Update the status of companies"
}

func (c *UpdateCompanyStatusCapability) NewInput() UpdateCompanyStatusInput {
	return UpdateCompanyStatusInput{}
}

func (c *UpdateCompanyStatusCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *UpdateCompanyStatusCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *UpdateCompanyStatusCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

func (c *UpdateCompanyStatusCapability) ValidateInput(input UpdateCompanyStatusInput) error {
	if input.OrganizationID == "" {
		err := errors.New("required parameter missing: OrganizationID")
		return err
	}

	if input.IcpFit != enum.IcpNotSet && len(input.IcpFitRationale) != 3 {
		err := errors.New("ICP fit data not complete")
		return err
	}

	if input.IcpFit == enum.IcpNotSet && input.OrganizationRelationship == "" && input.OrganizationStage == "" {
		err := errors.New("There's nothing to update")
		return err
	}
	return nil
}

func (c *UpdateCompanyStatusCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[UpdateCompanyStatusInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, NoOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateCompanyStatusCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "executionContainer", executionContainer)

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}

	var err error
	switch {
	case executionContainer.InputData.IcpFit == enum.IcpIsFit:
		err = c.processICPFit(ctx, executionContainer.InputData.OrganizationID, executionContainer.InputData.IcpFitRationale)
		if err != nil {
			tracing.TraceErr(span, err)
			return enum.CapabilityExecutionError, NoOutput{}, err
		}

	case executionContainer.InputData.IcpFit == enum.IcpNotFit:
		err = c.processICPNotAFit(ctx, executionContainer.InputData.OrganizationID, executionContainer.InputData.IcpFitRationale)
		if err != nil {
			tracing.TraceErr(span, err)
			return enum.CapabilityExecutionError, NoOutput{}, err
		}

	default:
		err = errors.New("Not implemented yet")
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, NoOutput{}, err
	}
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, NoOutput{}, err
	}

	publishErr := c.publishIcpFitEvent(ctx, executionContainer.InputData.IcpFit, executionContainer.AgentExecutionID)
	if publishErr != nil {
		tracing.TraceErr(span, publishErr)
		return enum.CapabilityExecutionError, NoOutput{}, err
	}

	return enum.CapabilityExecutionCompleted, NoOutput{}, nil
}

func (c *UpdateCompanyStatusCapability) processICPFit(ctx context.Context, organizationID string, reasons []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ICPQualificationCapability.processICPFit")
	defer span.Finish()
	tracing.TagComponentService(span)

	_, err := c.organizationService.Save(ctx, nil, &organizationID, data_fields.OrganizationFields{
		Relationship:  utils.ToPtr(neo4jenum.OrganizationRelationshipProspect),
		Stage:         utils.ToPtr(neo4jenum.Target),
		IcpFit:        utils.ToPtr(enum.IcpIsFit),
		IcpFitReasons: &reasons,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (c *UpdateCompanyStatusCapability) processICPNotAFit(ctx context.Context, organizationID string, reasons []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ICPQualificationCapability.processICPFit")
	defer span.Finish()
	tracing.TagComponentService(span)

	_, err := c.organizationService.Save(ctx, nil, &organizationID, data_fields.OrganizationFields{
		Relationship:  utils.ToPtr(neo4jenum.OrganizationRelationshipNotAFit),
		Stage:         utils.ToPtr(neo4jenum.Unqualified),
		IcpFit:        utils.ToPtr(enum.IcpNotFit),
		IcpFitReasons: &reasons,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (c *UpdateCompanyStatusCapability) publishIcpFitEvent(ctx context.Context, icpFitResult enum.IcpFit, agentExecutionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateCompanyStatusCapability.publishIcpFitEvent")
	defer span.Finish()
	tracing.TagComponentService(span)

	switch icpFitResult {
	case enum.IcpIsFit:
		return c.events.Publisher.PublishFanoutEvent(ctx, agentExecutionID, model.AGENT_EXECUTION, dto.IcpFit{
			AgentExecutionId: agentExecutionID,
		})

	case enum.IcpNotFit:
		return c.events.Publisher.PublishFanoutEvent(ctx, agentExecutionID, model.AGENT_EXECUTION, dto.IcpNotAFit{
			AgentExecutionId: agentExecutionID,
		})

	default:
		return errors.New("ICP Fit not set")
	}
}
