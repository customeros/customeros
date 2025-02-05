package agent_capability

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type UpdateCompanyStatusCapability struct{}

type UpdateCompanyStatusInput struct{}

type UpdateCompanyStatusOutput struct {
	IcpFit          enum.IcpFit `json:"icpFit"`
	IcpFitRationale []string    `json:"icpFitRationale"`
}

func NewUpdateCompanyStatusCapability() *UpdateCompanyStatusCapability {
	return &UpdateCompanyStatusCapability{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[UpdateCompanyStatusInput, UpdateCompanyStatusOutput, NoConfig] = (*UpdateCompanyStatusCapability)(nil)
)

func (c *UpdateCompanyStatusCapability) Type() enum.AgentCapability {
	return enum.CapabilityEvaluateCompanyICPFit
}

func (c *UpdateCompanyStatusCapability) Name() string {
	return "Evaluate company for ICP fit"
}

func (c *UpdateCompanyStatusCapability) NewInput() EvaluateICPFitInput {
	return EvaluateICPFitInput{}
}

func (c *UpdateCompanyStatusCapability) NewConfig() EvaluateICPFitConfig {
	return EvaluateICPFitConfig{}
}

func (c *UpdateCompanyStatusCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *UpdateCompanyStatusCapability) ValidateConfig(config EvaluateICPFitConfig) error {
	if config.QualificationCriteria.Value == "" {
		return errors.New("missing required config: QualificationCriteria")
	}
	return nil
}

func (c *UpdateCompanyStatusCapability) ValidateInput(input EvaluateICPFitInput) error {
	if input.OrganizationID == "" {
		return errors.New("missing required input: OrganizationID")
	}
	return nil
}

func (c *UpdateCompanyStatusCapability) Execute(ctx context.Context, data EvaluateICPFitInput, config EvaluateICPFitConfig) (EvaluateICPFitOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateCompanyStatusCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", data)
	tracing.LogObjectAsJson(span, "config", config)
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
