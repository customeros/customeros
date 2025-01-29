package agent_capability

import (
	"context"
	"encoding/json"
	"fmt"

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
	markdownEventService interfaces.MarkdownEventService
	organizationService  interfaces.OrganizationService
}

type ICPQualificationInput struct {
	OrganizationID string `json:"organizationId"`
}

type ICPQualificationOutput struct {
	CapabilityOutput
	IsICPFit        string   `json:"isIcpFit"`
	IcpFitRationale []string `json:"icpFitRationale"`
}

type ICPFit string

const (
	ICPIsFit   ICPFit = "icp_fit"
	ICPNotAFit ICPFit = "icp_not_a_fit"
	ICPNotSet  ICPFit = "not_set"
)

type ICPAnswer struct {
	ICPFit  bool     `json:"icp_fit"`
	Reasons []string `json:"reasons"`
}

func (c *ICPQualificationCapability) ValidateConfig(config NoConfig) error {
	// if config.QualificationCriteria == "" {
	// 	return errors.New("missing required input: QualificationCriteria")
	// }
	// if config.DisqualificationCriteria == "" {
	// 	return errors.New("missing required input: DisqualificationCriteria")
	// }
	return nil
}

func (c *ICPQualificationCapability) ValidateInput(input ICPQualificationInput) error {
	if input.OrganizationID == "" {
		return errors.New("missing required input: OrganizationID")
	}
	return nil
}

func (c *ICPQualificationCapability) GetInput() any {
	return &ICPQualificationInput{}
}

func (c *ICPQualificationCapability) GetConfig() any {
	return &NoConfig{}
}

func (c *ICPQualificationCapability) GetOutput() any {
	return &ICPQualificationOutput{}
}

func NewICPQualificationCapability(postgres *postgres_repository.Repositories, aiService interfaces.AIService) *ICPQualificationCapability {
	return &ICPQualificationCapability{
		postgresRepositories: postgres,
		aiService:            aiService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapabilityExecution[ICPQualificationInput, ICPQualificationOutput, NoConfig] = (*ICPQualificationCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                            = (*ICPQualificationCapability)(nil)
)

func (c *ICPQualificationCapability) Execute(ctx context.Context, data ICPQualificationInput, config NoConfig) (ICPQualificationOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.executeICPQualification")
	defer span.Finish()
	tracing.TagComponentService(span)

	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)

	result := ICPQualificationOutput{
		IsICPFit: string(ICPNotSet),
	}

	err := c.runExecutionValidation(ctx, data, config)
	if err != nil {
		tracing.TraceErr(span, err)
		result.CapabilityOutput.ExecutionValidated = false
		return result, err
	}
	result.CapabilityOutput.ExecutionValidated = true

	// get ICP definition
	icpQualification := ""
	icpDisqualification := ""

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
	systemPrompt, content := c.buildPrompts(icpQualification, icpDisqualification, company)

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
		result.CapabilityOutput.Completed = false
		return result, err
	}
	result.IcpFitRationale = parsedAnswer.Reasons

	// save to timeline
	timelineContent := c.buildTimelineEvent(parsedAnswer)
	source := enum.SourceAgent
	_, err = c.markdownEventService.Save(ctx, nil, nil, data_fields.MarkdownEventFields{
		Source:         &source,
		OrganizationId: &data.OrganizationID,
		Content:        &timelineContent,
		CreatedAt:      utils.NowPtr(),
	})

	// disqualify as lead if not a fit
	if parsedAnswer.ICPFit {
		result.IsICPFit = string(ICPIsFit)
		err := c.processICPFit(ctx, data.OrganizationID)
		if err != nil {
			tracing.TraceErr(span, err)
			result.CapabilityOutput.Completed = false
			return result, err
		}
	} else {
		result.IsICPFit = string(ICPNotAFit)
		err := c.processICPNotAFit(ctx, data.OrganizationID)
		if err != nil {
			tracing.TraceErr(span, err)
			result.CapabilityOutput.Completed = false
			return result, err
		}
	}

	result.CapabilityOutput.Completed = true
	return result, nil
}

func (c *ICPQualificationCapability) processICPFit(ctx context.Context, organizationID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ICPQualificationCapability.processICPFit")
	defer span.Finish()
	tracing.TagComponentService(span)

	relationship := neo4jenum.OrganizationRelationshipProspect
	stage := neo4jenum.Target
	_, err := c.organizationService.Save(ctx, nil, &organizationID, data_fields.OrganizationFields{
		Relationship: &relationship,
		Stage:        &stage,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (c *ICPQualificationCapability) processICPNotAFit(ctx context.Context, organizationID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ICPQualificationCapability.processICPFit")
	defer span.Finish()
	tracing.TagComponentService(span)

	relationship := neo4jenum.OrganizationRelationshipNotAFit
	stage := neo4jenum.Unqualified
	_, err := c.organizationService.Save(ctx, nil, &organizationID, data_fields.OrganizationFields{
		Relationship: &relationship,
		Stage:        &stage,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (c *ICPQualificationCapability) buildTimelineEvent(answer *ICPAnswer) string {
	md := fmt.Sprintf("**ICP Fit:** %v\n", answer.ICPFit)

	// Add reasons if there are any
	if len(answer.Reasons) > 0 {
		md += "**Reasons:**\n"
		for _, reason := range answer.Reasons {
			md += fmt.Sprintf("* %s\n", reason)
		}
	}

	return md
}

func (c ICPQualificationCapability) buildPrompts(icpQualification, icpDisqualification string, company *postgres_entity.GlobalOrganization) (string, string) {
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
        Company Description 1: %s
        Company Description 2: %s
        Company Description 3: %s
        Company Description 4: %s
        Company Description 5: %s
        `, icpQualification, icpDisqualification,
		company.Name, company.PrimaryDomain, company.YearFounded, company.EmployeeCount,
		company.City, company.Region, company.CountryA2, company.IndustryNaicsCode,
		company.IndustryNaicsName, company.Description, company.SourceDescription1,
		company.SourceDescription2, company.SourceDescription3, company.SourceDescription4)

	return systemPrompt, content
}

func (c *ICPQualificationCapability) runExecutionValidation(ctx context.Context, data ICPQualificationInput, config NoConfig) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ICPQualificationCapability.runValidation")
	defer span.Finish()
	tracing.TagComponentService(span)

	if err := c.ValidateInput(data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return err
	}
	if err := c.ValidateConfig(config); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return err
	}
	return nil
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

// ExecuteUntyped implements the AgentCapabilityUntyped interface.
// It casts the generic input and config to the specific types and delegates to the typed Execute method.
func (c *ICPQualificationCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
	typedInput, ok := input.(*ICPQualificationInput)
	if !ok || typedInput == nil {
		return nil, fmt.Errorf("invalid input type: expected ICPQualificationInput")
	}

	typedConfig, ok := config.(*NoConfig)
	if !ok || typedConfig == nil {
		return nil, fmt.Errorf("invalid config type: expected NoCOnfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
