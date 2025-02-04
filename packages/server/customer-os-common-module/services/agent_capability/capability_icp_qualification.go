package agent_capability

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type ICPQualificationCapability struct {
	postgresRepositories *postgres_repository.Repositories
	aiService            interfaces.AIService
	organizationService  interfaces.OrganizationService
}

type ICPQualificationInput struct {
	OrganizationID string `json:"organizationId"`
}

type ICPQualificationOutput struct {
	CapabilityOutput
	IcpFit          enum.IcpFit `json:"icpFit"`
	IcpFitRationale []string    `json:"icpFitRationale"`
}

type ICPQualificationConfig struct {
	QualificationCriteria    QualificationCriteriaConfig    `json:"qualificationCriteria"`
	DisqualificationCriteria DisqualificationCriteriaConfig `json:"disqualificationCriteria"`
}

type QualificationCriteriaConfig struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

type DisqualificationCriteriaConfig struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

func (c *ICPQualificationConfig) Validate() bool {
	isValid := true

	// validate qualification criteria is not empty
	if c.QualificationCriteria.Value == "" {
		c.QualificationCriteria.Error = "Please provide qualification criteria."
		isValid = false
	} else {
		c.QualificationCriteria.Error = ""
	}

	return isValid
}

func NewICPQualificationCapability(postgres *postgres_repository.Repositories, aiService interfaces.AIService, organizationServices interfaces.OrganizationService) *ICPQualificationCapability {
	return &ICPQualificationCapability{
		postgresRepositories: postgres,
		aiService:            aiService,
		organizationService:  organizationServices,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[ICPQualificationInput, ICPQualificationOutput, ICPQualificationConfig] = (*ICPQualificationCapability)(nil)
)

type ICPAnswer struct {
	ICPFit  bool     `json:"icp_fit"`
	Reasons []string `json:"reasons"`
}

func (c *ICPQualificationCapability) Type() enum.AgentCapability {
	return enum.CapabilityEvaluateCompanyICPFit
}

func (c *ICPQualificationCapability) Name() string {
	return "Evaluate company for ICP fit"
}

func (c *ICPQualificationCapability) NewInput() ICPQualificationInput {
	return ICPQualificationInput{}
}

func (c *ICPQualificationCapability) NewConfig() ICPQualificationConfig {
	return ICPQualificationConfig{}
}

func (c *ICPQualificationCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *ICPQualificationCapability) ValidateConfig(config ICPQualificationConfig) error {
	if config.QualificationCriteria.Value == "" {
		return errors.New("missing required config: QualificationCriteria")
	}
	return nil
}

func (c *ICPQualificationCapability) ValidateInput(input ICPQualificationInput) error {
	if input.OrganizationID == "" {
		return errors.New("missing required input: OrganizationID")
	}
	return nil
}

func (c *ICPQualificationCapability) Execute(ctx context.Context, data ICPQualificationInput, config ICPQualificationConfig) (ICPQualificationOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ICPQualificationCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", data)
	tracing.LogObjectAsJson(span, "config", config)

	result := ICPQualificationOutput{
		IcpFit: enum.IcpNotSet,
	}

	if err := c.ValidateInput(data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return result, err
	}
	if err := c.ValidateConfig(config); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return result, err
	}

	result.ExecutionValidated = true

	// get all company context
	primaryDomain, err := c.organizationService.GetPrimaryDomainByOrgID(ctx, data.OrganizationID)
	if err != nil {
		tracing.TraceErr(span, err)
		result.CapabilityOutput.Completed = false
		return result, err
	}
	company, err := c.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		result.CapabilityOutput.Completed = false
		return result, err
	}
	if company == nil {
		err := errors.New("Company not set in global orgs, cannot run ICP qualification")
		tracing.TraceErr(span, err)
		result.CapabilityOutput.Completed = false
		return result, err
	}

	// build prompt
	systemPrompt, content := c.buildPrompts(config.QualificationCriteria.Value, config.DisqualificationCriteria.Value, company)

	// askAI
	answer, err := c.aiService.AskAI(ctx, enum.AIModelAnthropicHaiku, systemPrompt, content)
	if err != nil {
		tracing.TraceErr(span, err)
		result.CapabilityOutput.Completed = false
		return result, err
	}

	// parse answer
	parsedAnswer, err := c.parseAnswer(ctx, *answer)
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}
	result.IcpFitRationale = parsedAnswer.Reasons

	// disqualify as lead if not a fit
	if parsedAnswer.ICPFit {
		err = c.processICPFit(ctx, data.OrganizationID, parsedAnswer.Reasons)
		if err != nil {
			tracing.TraceErr(span, err)
			return result, err
		}
		result.IcpFit = enum.IcpIsFit
	} else {
		err = c.processICPNotAFit(ctx, data.OrganizationID, parsedAnswer.Reasons)
		if err != nil {
			tracing.TraceErr(span, err)
			return result, err
		}
		result.IcpFit = enum.IcpNotFit
	}
	result.IcpFitRationale = parsedAnswer.Reasons

	result.Completed = true
	tracing.LogObjectAsJson(span, "result", result)
	return result, nil
}

func (c *ICPQualificationCapability) processICPFit(ctx context.Context, organizationID string, reasons []string) error {
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

func (c *ICPQualificationCapability) processICPNotAFit(ctx context.Context, organizationID string, reasons []string) error {
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

func (c *ICPQualificationCapability) buildPrompts(icpQualification, icpDisqualification string, company *postgres_entity.GlobalOrganization) (string, string) {
	systemPrompt := `You are a world class company analyst. Your objective is to determine whether the company I provide you fits our ideal customer profile or not. I will provide you with three datasets: 1. a description of our ideal customer, 2. criteria that automatically disqualifies companies, and 3. all the context I have about the company, including their name, location, and several descriptions taken from their website and linkedin pages.

Analyze the company and respond in this exact JSON format:
{
    "icp_fit": false,
    "reasons": [
        "First key point about qualification",
        "Second key point",
        "Third key point"
    ]
}

Important: Always provide exactly three reasons, and format as valid JSON.`

	var descLines []string
	descriptions := []string{company.SourceDescription1, company.SourceDescription2, company.SourceDescription3, company.SourceDescription4, company.SourceDescription5}
	for i, d := range descriptions {
		if strings.TrimSpace(d) != "" {
			descLines = append(descLines, fmt.Sprintf("Description Line %d: %s", i+1, d))
		}
	}
	additionalCompanyDescriptions := strings.Join(descLines, ";")

	content := fmt.Sprintf(`
        ICP Qualificaton Criteria: %s
        ICP Disqualification Criteria: %s
        Company Name: %s
        Company Domain: %s
        Year Founded: %d
        Employee Count: %d
        Location: %s, %s, %s
        Industry NAICS Code: %s
        Industry Name: %s
		Company Description: %s
		Additional Company Descriptions: %s
        `, icpQualification, icpDisqualification,
		company.Name, company.PrimaryDomain, company.YearFounded, company.EmployeeCount,
		company.City, company.Region, company.CountryA2, company.IndustryNaicsCode,
		company.IndustryNaicsName, company.Description, additionalCompanyDescriptions)

	return systemPrompt, content
}

func (c *ICPQualificationCapability) parseAnswer(ctx context.Context, answer string) (*ICPAnswer, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ICPQualificationCapability.parseAnswer")
	defer span.Finish()
	tracing.TagComponentService(span)

	var result ICPAnswer
	err := json.Unmarshal([]byte(answer), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate we got exactly 3 reasons
	if len(result.Reasons) != 3 {
		return nil, fmt.Errorf("expected exactly 3 reasons, got %d", len(result.Reasons))
	}

	return &result, nil
}
