package agent_capability

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

const MinICPCompanyExamples = 5

type EvaluateICPFitCapability struct {
	aiService         interfaces.AIService
	webscraperService interfaces.WebscraperService
}

type EvaluateICPFitInput struct {
	OrganizationID     string `json:"organizationId"`
	CompanyName        string `json:"companyName"`
	PrimaryDomain      string `json:"primaryDomain"`
	IndustryNAICSName  string `json:"industryName"`
	YearCompanyFounded string `json:"yearCompanyFounded"`
	EmployeeCount      int64  `json:"employeeCount"`
	CompanyCity        string `json:"companyCity"`
	CompanyRegion      string `json:"companyRegion"`
	CompanyCountryA2   string `json:"companyCountry"`
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

	if len(c.ICPCompanyExamples.Value) < MinICPCompanyExamples {
		c.ICPCompanyExamples.Error = "Add at least 5 websites"
		isValid = false
	} else {
		c.ICPCompanyExamples.Error = ""
	}

	return isValid
}

func NewEvaluateICPFitCapability(aiService interfaces.AIService, webscraperService interfaces.WebscraperService) *EvaluateICPFitCapability {
	return &EvaluateICPFitCapability{
		aiService:         aiService,
		webscraperService: webscraperService,
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
	return "Qualify if companies fit your ICP"
}

func (c *EvaluateICPFitCapability) NewInput() EvaluateICPFitInput {
	return EvaluateICPFitInput{}
}

func (c *EvaluateICPFitCapability) NewConfig() EvaluateICPFitConfig {
	return EvaluateICPFitConfig{}
}

func (c *EvaluateICPFitCapability) DefaultConfig() any {
	config := c.NewConfig()
	config.ICPCompanyExamples.Value = []string{}
	return &config
}

func (c *EvaluateICPFitCapability) DefaultActive() bool {
	return true
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

func (c *EvaluateICPFitCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[EvaluateICPFitInput, EvaluateICPFitConfig]) (enum.CapabilityExecutionStatus, EvaluateICPFitOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EvaluateICPFitCapability.Execute")
	defer spans.Finish()

	spans.LogObjectAsJson("executionContainer", executionContainer)

	result := EvaluateICPFitOutput{
		IcpFit: enum.IcpNotSet,
	}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	// get ICP profile

	// build prompt
	systemPrompt, content, err := c.buildPrompts(ctx, executionContainer)
	if err != nil {
		spans.TraceError(err)
		return enum.CapabilityExecutionError, result, err
	}

	// askAI
	answer, err := c.aiService.AskAI(ctx, interfaces.AskAIRequest{
		Model:        enum.AIModelGemini,
		SystemPrompt: &systemPrompt,
		Prompt:       &content,
		OutputFormat: enum.AIOutputJson,
	})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to ask AI"))
		if answer != nil {
			spans.LogKV("result.answer", *answer)
		} else {
			spans.LogKV("result.answer", "nil")
		}
		return enum.CapabilityExecutionError, result, err
	}

	// parse answer
	parsedAnswer, err := c.parseAnswer(ctx, *answer)
	if err != nil {
		spans.TraceError(err)
		return enum.CapabilityExecutionError, result, err
	}

	result.IcpFitRationale = parsedAnswer.Reasons
	if parsedAnswer.ICPFit {
		result.IcpFit = enum.IcpIsFit
	} else {
		result.IcpFit = enum.IcpNotFit
	}
	spans.LogObjectAsJson("result", result)

	return enum.CapabilityExecutionCompleted, result, nil
}

func (c *EvaluateICPFitCapability) buildPrompts(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[EvaluateICPFitInput, EvaluateICPFitConfig]) (string, string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EvaluateICPFitCapability.buildPrompts")
	defer spans.Finish()

	company := executionContainer.InputData

	homepage, err := c.webscraperService.Scrape(ctx, company.PrimaryDomain)
	if err != nil {
		spans.LogKV("scapeResults", "none")
		homepage = ""
	}

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

Important:
- Always output exactly three reasons.
- Your entire response must be valid JSON and nothing else.
- Follow the JSON format exactly.
- Pay special attention to location criteria in your decision making.`

	var contentBuilder strings.Builder

	contentBuilder.WriteString("\n        ICP Qualificaton Criteria: ")
	contentBuilder.WriteString(executionContainer.ConfigData.QualificationCriteria.Value)
	contentBuilder.WriteString("\n        ICP Disqualification Criteria: ")
	contentBuilder.WriteString(executionContainer.ConfigData.DisqualificationCriteria.Value)
	contentBuilder.WriteString("\n        Company Name: ")
	contentBuilder.WriteString(company.CompanyName)
	contentBuilder.WriteString("\n        Company Domain: ")
	contentBuilder.WriteString(company.PrimaryDomain)
	contentBuilder.WriteString("\n        Year Founded: ")
	contentBuilder.WriteString(company.YearCompanyFounded)
	contentBuilder.WriteString("\n        Employee Count: ")
	contentBuilder.WriteString(strconv.Itoa(int(company.EmployeeCount)))
	contentBuilder.WriteString("\n        Location: ")
	contentBuilder.WriteString(company.CompanyCity)
	contentBuilder.WriteString(", ")
	contentBuilder.WriteString(company.CompanyRegion)
	contentBuilder.WriteString(", ")
	contentBuilder.WriteString(company.CompanyCountryA2)
	contentBuilder.WriteString("\n        Industry Name: ")
	contentBuilder.WriteString(company.IndustryNAICSName)
	if homepage != "" {
		contentBuilder.WriteString("\n\t\tCompany Homepage: ")
		contentBuilder.WriteString(homepage)
	}

	content := contentBuilder.String()

	return systemPrompt, content, nil
}

func (c *EvaluateICPFitCapability) parseAnswer(ctx context.Context, answer string) (*ICPAnswer, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EvaluateICPFitCapability.parseAnswer")
	defer spans.Finish()

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
