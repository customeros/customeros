package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type UpdateCompanyStatusCapability struct {
	postgres            *postgres_repository.Repositories
	events              *events.EventsService
	organizationService interfaces.OrganizationService
}

type UpdateCompanyStatusInput struct {
	OrganizationID           string      `json:"organizationId"`
	OrganizationRelationship string      `json:"organizationRelationship"`
	OrganizationStage        string      `json:"organizationStage"`
	IcpFit                   enum.IcpFit `json:"icpFit"`
	IcpFitRationale          []string    `json:"icpFitRationale"`
}

func NewUpdateCompanyStatusCapability(postgres *postgres_repository.Repositories, events *events.EventsService, orgSrv interfaces.OrganizationService) *UpdateCompanyStatusCapability {
	return &UpdateCompanyStatusCapability{
		postgres:            postgres,
		events:              events,
		organizationService: orgSrv,
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

func (c *UpdateCompanyStatusCapability) DefaultActive() bool {
	return true
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
	spans, ctx := telemetry.StartServiceSpan(ctx, "UpdateCompanyStatusCapability.Execute")
	defer spans.Finish()

	spans.LogObjectAsJson("executionContainer", executionContainer)

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}

	var err error
	switch {
	case executionContainer.InputData.IcpFit == enum.IcpIsFit:
		err = c.processICPFit(ctx, executionContainer.InputData.OrganizationID, executionContainer.InputData.IcpFitRationale)
		if err != nil {
			spans.TraceError(err)
			return enum.CapabilityExecutionError, NoOutput{}, err
		}

	case executionContainer.InputData.IcpFit == enum.IcpNotFit:
		err = c.processICPNotAFit(ctx, executionContainer.InputData.OrganizationID, executionContainer.InputData.IcpFitRationale)
		if err != nil {
			spans.TraceError(err)
			return enum.CapabilityExecutionError, NoOutput{}, err
		}

	default:
		err = errors.New("Not implemented yet")
		spans.TraceError(err)
		return enum.CapabilityExecutionError, NoOutput{}, err
	}
	if err != nil {
		spans.TraceError(err)
		return enum.CapabilityExecutionError, NoOutput{}, err
	}

	err = c.postgres.AgentExecutionRepository.GoalAchieved(ctx, executionContainer.AgentExecutionID, true, nil)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to publish ignore email event"))
	}

	return enum.CapabilityExecutionCompleted, NoOutput{}, nil
}

func (c *UpdateCompanyStatusCapability) processICPFit(ctx context.Context, organizationID string, reasons []string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UpdateCompanyStatusCapability.processICPFit")
	defer spans.Finish()

	spans.TagEntity(organizationID)

	_, err := c.organizationService.Save(ctx, nil, &organizationID, data_fields.OrganizationFields{
		Relationship:  utils.ToPtr(neo4jenum.OrganizationRelationshipProspect),
		Stage:         utils.ToPtr(enum.Target),
		IcpFit:        utils.ToPtr(enum.IcpIsFit),
		IcpFitReasons: &reasons,
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}

func (c *UpdateCompanyStatusCapability) processICPNotAFit(ctx context.Context, organizationID string, reasons []string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "UpdateCompanyStatusCapability.processICPNotAFit")
	defer spans.Finish()

	spans.TagEntity(organizationID)

	_, err := c.organizationService.Save(ctx, nil, &organizationID, data_fields.OrganizationFields{
		Relationship:  utils.ToPtr(neo4jenum.OrganizationRelationshipNotAFit),
		Stage:         utils.ToPtr(enum.Unqualified),
		IcpFit:        utils.ToPtr(enum.IcpNotFit),
		IcpFitReasons: &reasons,
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}
