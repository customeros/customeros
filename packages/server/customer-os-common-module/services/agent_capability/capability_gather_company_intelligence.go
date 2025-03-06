package agent_capability

import (
	"context"
	"github.com/opentracing/opentracing-go/log"
	"strconv"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type GatherCompanyIntelligenceCapability struct {
	postgresRepositories *postgres_repository.Repositories
	organizationService  interfaces.OrganizationService
}

type GatherCompanyIntilligenceInput struct {
	Domain         string `json:"domain"`
	OrganizationID string `json:"organizationId"`
}

type GatherCompanyIntelligenceOutput struct {
	CompanyName         string              `json:"companyName"`
	PrimaryDomain       string              `json:"primaryDomain"`
	CompanyDescriptions CompanyDescriptions `json:"companyDescriptions"`
	IndustryNAICSName   string              `json:"industryName"`
	YearCompanyFounded  string              `json:"yearCompanyFounded"`
	EmployeeCount       int64               `json:"employeeCount"`
	CompanyCity         string              `json:"companyCity"`
	CompanyRegion       string              `json:"companyRegion"`
	CompanyCountryA2    string              `json:"companyCountry"`
}

type CompanyDescriptions struct {
	Description  string `json:"description"`
	Description2 string `json:"description2"`
	Description3 string `json:"description3"`
	Description4 string `json:"description4"`
	Description5 string `json:"description5"`
}

func NewGatherCompanyIntelligenceCapability(
	postgres *postgres_repository.Repositories,
	organizationService interfaces.OrganizationService,
) *GatherCompanyIntelligenceCapability {
	return &GatherCompanyIntelligenceCapability{
		postgresRepositories: postgres,
		organizationService:  organizationService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[GatherCompanyIntilligenceInput, GatherCompanyIntelligenceOutput, postgres_entity.NoConfig] = (*GatherCompanyIntelligenceCapability)(nil)
)

func (c *GatherCompanyIntelligenceCapability) Type() enum.AgentCapability {
	return enum.CapabilityGatherCompanyIntelligence
}

func (c *GatherCompanyIntelligenceCapability) Name() string {
	return "Gather intelligence on companies"
}

func (c *GatherCompanyIntelligenceCapability) NewInput() GatherCompanyIntilligenceInput {
	return GatherCompanyIntilligenceInput{}
}

func (c *GatherCompanyIntelligenceCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *GatherCompanyIntelligenceCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *GatherCompanyIntelligenceCapability) DefaultActive() bool {
	return true
}

func (c *GatherCompanyIntelligenceCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

func (c *GatherCompanyIntelligenceCapability) ValidateInput(input GatherCompanyIntilligenceInput) error {
	if input.OrganizationID == "" && input.Domain == "" {
		return errors.New("missing required input: OrganizationID or Domain")
	}
	return nil
}

func (c *GatherCompanyIntelligenceCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[GatherCompanyIntilligenceInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, GatherCompanyIntelligenceOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GatherCompanyIntelligenceCapability.Execute")
	defer span.Finish()
	tracing.SetDefaultAgentCapabilitySpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "executionContainer", executionContainer)

	result := GatherCompanyIntelligenceOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	// get all company context
	primaryDomain, err := c.organizationService.GetPrimaryDomainByOrgID(ctx, executionContainer.InputData.OrganizationID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}
	company, err := c.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}
	if company == nil {
		span.LogFields(log.String("result", "Company not set in global orgs"))
		return enum.CapabilityExecutionRetry, result, nil
	}

	result.CompanyName = company.Name
	result.PrimaryDomain = company.PrimaryDomain
	result.CompanyDescriptions.Description = company.Description
	result.CompanyDescriptions.Description2 = company.SourceDescription1
	result.CompanyDescriptions.Description3 = company.SourceDescription2
	result.CompanyDescriptions.Description4 = company.SourceDescription3
	result.CompanyDescriptions.Description5 = company.SourceDescription4
	result.IndustryNAICSName = company.IndustryNaicsName
	result.YearCompanyFounded = strconv.Itoa(company.YearFounded)
	result.EmployeeCount = company.EmployeeCount
	result.CompanyCity = company.City
	result.CompanyRegion = company.Region
	result.CompanyCountryA2 = company.CountryA2
	return enum.CapabilityExecutionCompleted, result, nil
}
