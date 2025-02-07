package agent_capability

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

const MinICPCompanyExamples = 5

type EvaluateICPFitCapability struct {
	aiService interfaces.AIService
	events    *events.EventsService
}

type EvaluateICPFitInput struct {
	OrganizationID      string              `json:"organizationId"`
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

type EvaluateICPFitOutput struct {
	IcpFit          enum.IcpFit `json:"icpFit"`
	IcpFitRationale []string    `json:"icpFitRationale"`
}

type EvaluateICPFitConfig struct {
	ICPCompanyExamples       ConfigMultipleValues `json:"icpCompanyExamples"`
	QualificationCriteria    ConfigSingleValue    `json:"qualificationCriteria"`
	DisqualificationCriteria ConfigSingleValue    `json:"disqualificationCriteria"`
}

func (c *EvaluateICPFitConfig) Validate() bool {
	isValid := true

	if c.QualificationCriteria.Value == "" {
		c.QualificationCriteria.Error = "Please provide a qualification criteria"
		isValid = false
	} else {
		c.QualificationCriteria.Error = ""
	}
	if len(c.ICPCompanyExamples.Value) < MinICPCompanyExamples {
		c.ICPCompanyExamples.Error = "Please provide at least 5 companies that match your ICP."
		isValid = false
	} else {
		c.ICPCompanyExamples.Error = ""
	}

	return isValid
}

func NewEvaluateICPFitCapability(aiService interfaces.AIService, events *events.EventsService) *EvaluateICPFitCapability {
	return &EvaluateICPFitCapability{
		aiService: aiService,
		events:    events,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[EvaluateICPFitInput, EvaluateICPFitOutput, EvaluateICPFitConfig] = (*EvaluateICPFitCapability)(nil)
)

func (c *EvaluateICPFitCapability) Type() enum.AgentCapability {
	return enum.CapabilityEvaluateCompanyICPFit
}

func (c *EvaluateICPFitCapability) Name() string {
	return "Evaluate company for ICP fit"
}

func (c *EvaluateICPFitCapability) NewInput() EvaluateICPFitInput {
	return EvaluateICPFitInput{}
}

func (c *EvaluateICPFitCapability) NewConfig() EvaluateICPFitConfig {
	return EvaluateICPFitConfig{}
}

func (c *EvaluateICPFitCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *EvaluateICPFitCapability) ValidateConfig(config EvaluateICPFitConfig) error {
	if len(config.ICPCompanyExamples.Value) < MinICPCompanyExamples {
		return errors.New("missing required config: ICPCompanyExamples")
	}
	return nil
}

func (c *EvaluateICPFitCapability) ValidateInput(input EvaluateICPFitInput) error {
	switch {
	case input.OrganizationID == "":
		return errors.New("missing required input: OrganizationID")
	case input.PrimaryDomain == "":
		return errors.New("missing required input: PrimaryDomain")
	case input.CompanyName == "":
		return errors.New("missing required input: CompanyName")
	case input.CompanyDescriptions.Description1 == "" && input.CompanyDescriptions.Description2 == "":
		return errors.New("missing required input: CompanyDescription")
	case input.IndustryNAICSName == "":
		return errors.New("missing required input: IndustryNAICSName")
	case input.EmployeeCount == 0:
		return errors.New("missing required input: EmployeeCount cannot be 0")
	case input.CompanyCountryA2 == "":
		return errors.New("missing required input: CompanyCountryA2")
	default:
		return nil
	}
}

type ICPAnswer struct {
	ICPFit  bool     `json:"icp_fit"`
	Reasons []string `json:"reasons"`
}

func (c *EvaluateICPFitCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[EvaluateICPFitInput, EvaluateICPFitConfig]) (bool, EvaluateICPFitOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EvaluateICPFitCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "executionContainer", executionContainer)

	result := EvaluateICPFitOutput{
		IcpFit: enum.IcpNotSet,
	}

	if executionContainer.InputData.EmployeeCount == 0 {
		result.IcpFit = enum.IcpNotFit
		return true, result, nil
	}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return false, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return false, result, err
	}

	// build prompt
	systemPrompt, content := c.buildPrompts(executionContainer)

	// askAI
	answer, err := c.aiService.AskAI(ctx, enum.AIModelAnthropicHaiku, systemPrompt, content)
	if err != nil {
		tracing.TraceErr(span, err)
		return true, result, err
	}

	// parse answer
	parsedAnswer, err := c.parseAnswer(ctx, *answer)
	if err != nil {
		tracing.TraceErr(span, err)
		return true, result, err
	}
	result.IcpFitRationale = parsedAnswer.Reasons
	tracing.LogObjectAsJson(span, "result", result)

	err = c.publishIcpFitEvent(ctx, result.IcpFit, executionContainer.AgentExecutionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return true, result, err
	}

	return true, result, nil
}

func (c *EvaluateICPFitCapability) publishIcpFitEvent(ctx context.Context, icpFitResult enum.IcpFit, agentExecutionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EvaluateICPFitCapability.publishIcpFitEvent")
	defer span.Finish()
	tracing.TagComponentService(span)

	switch icpFitResult {
	case enum.IcpIsFit:
		return c.events.Publisher.PublishFanoutEvent(ctx, agentExecutionID, model.AGENT_EXECUTION, dto.IcpFit{})

	case enum.IcpNotFit:
		return c.events.Publisher.PublishFanoutEvent(ctx, agentExecutionID, model.AGENT_EXECUTION, dto.IcpNotAFit{})

	default:
		return errors.New("ICP Fit not set")
	}
}

func (c *EvaluateICPFitCapability) buildPrompts(executionContainer interfaces.TypedExecutionContainer[EvaluateICPFitInput, EvaluateICPFitConfig]) (string, string) {
	company := executionContainer.InputData

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

Important: Always provide exactly three reasons, and format as valid JSON.  Please pay special attention to location criteria in your decision making.`

	var descLines []string
	descriptions := []string{company.CompanyDescriptions.Description2, company.CompanyDescriptions.Description3, company.CompanyDescriptions.Description4, company.CompanyDescriptions.Description5}
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
        Year Founded: %s
        Employee Count: %d
        Location: %s, %s, %s
        Industry Name: %s
		Company Description: %s
		Additional Company Descriptions: %s
        `, executionContainer.ConfigData.QualificationCriteria.Value, executionContainer.ConfigData.DisqualificationCriteria.Value,
		company.CompanyName, company.PrimaryDomain, company.YearCompanyFounded, company.EmployeeCount,
		company.CompanyCity, company.CompanyRegion, company.CompanyCountryA2,
		company.IndustryNAICSName, company.CompanyDescriptions.Description1, additionalCompanyDescriptions)

	return systemPrompt, content
}

func (c *EvaluateICPFitCapability) parseAnswer(ctx context.Context, answer string) (*ICPAnswer, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EvaluateICPFitCapability.parseAnswer")
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
