package agent_capability

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type EvaluateICPFitCapability struct {
	aiService interfaces.AIService
}

type EvaluateICPFitInput struct {
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

func (c *EvaluateICPFitConfig) Validate() bool {
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

func NewEvaluateICPFitCapability(aiService interfaces.AIService) *EvaluateICPFitCapability {
	return &EvaluateICPFitCapability{
		aiService: aiService,
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
	if config.QualificationCriteria.Value == "" {
		return errors.New("missing required config: QualificationCriteria")
	}
	return nil
}

func (c *EvaluateICPFitCapability) ValidateInput(input EvaluateICPFitInput) error {
	switch {
	case input.PrimaryDomain == "":
		return errors.New("missing required input: PrimaryDomain")
	case input.CompanyName == "":
		return errors.New("missing required input: CompanyName")
	case input.CompanyDescriptions.Description1 == "":
		return errors.New("missing required input: CompanyDescription")
	case input.IndustryNAICSName == "":
		return errors.New("missing required input: IndustryNAICSName")
	case input.EmployeeCount == 0:
		return errors.New("missing required input: EmployeeCount cannot be 0")
	case input.CompanyCountryA2 == "":
		return errors.New("missing required input: CompanyCountryA2")
	case input.CompanyName == "":
		return errors.New("missing required input: CompanyName")
	default:
		return nil
	}
}

type ICPAnswer struct {
	ICPFit  bool     `json:"icp_fit"`
	Reasons []string `json:"reasons"`
}

func (c *EvaluateICPFitCapability) Execute(ctx context.Context, data EvaluateICPFitInput, config EvaluateICPFitConfig) (EvaluateICPFitOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EvaluateICPFitCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", data)
	tracing.LogObjectAsJson(span, "config", config)

	result := EvaluateICPFitOutput{
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

	// build prompt
	systemPrompt, content := c.buildPrompts(config.QualificationCriteria.Value, config.DisqualificationCriteria.Value, data)

	// askAI
	answer, err := c.aiService.AskAI(ctx, enum.AIModelAnthropicHaiku, systemPrompt, content)
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}

	// parse answer
	parsedAnswer, err := c.parseAnswer(ctx, *answer)
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}
	result.IcpFitRationale = parsedAnswer.Reasons

	tracing.LogObjectAsJson(span, "result", result)
	return result, nil
}

func (c *EvaluateICPFitCapability) buildPrompts(icpQualification, icpDisqualification string, company EvaluateICPFitInput) (string, string) {
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
	descriptions := []string{company.CompanyDescriptions.Description2, company.CompanyDescriptions.Description3, company.CompanyDescriptions.Description4}
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
        Industry Name: %s
		Company Description: %s
		Additional Company Descriptions: %s
        `, icpQualification, icpDisqualification,
		company.CompanyName, company.PrimaryDomain, company.YearCompanyFounded, company.EmployeeCount,
		company.CompanyCity, company.CompanyRegion, company.CompanyCountryA2,
		company.IndustryNAICSName, company.CompanyDescriptions.Description1, additionalCompanyDescriptions)

	return systemPrompt, content
}

func (c *EvaluateICPFitCapability) parseAnswer(ctx context.Context, answer string) (*ICPAnswer, error) {
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
