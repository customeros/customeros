package agent_capability

import (
	"context"
	"fmt"
	"strings"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type ICPQualificationCapability struct {
	postgresRepositories *postgres_repository.Repositories
	aiService            interfaces.AIService
}

type ICPQualificationInput struct {
	QualificationCriteria    string `json:"organizationQualificationCriteria"`
	DisqualificationCriteria string `json:"organizationDisqualificationCriteria"`
	PrimaryDomain            string `json:"organizationPrimaryDomain"`
}

type ICPQualificationResult struct {
	isICPFit        string `json:"isIcpFit"`
	icpFitRationale string `json:"icpFitRationale"`
}

func (c *ICPQualificationCapability) ValidateConfig(config NoConfig) error {
	// TODO implement me
	panic("implement me")
}

func (c *ICPQualificationCapability) ValidateInput(input ICPQualificationInput) error {
	if input.PrimaryDomain == "" {
		return errors.New("missing required input: PromaryDomain")
	}
	if input.QualificationCriteria == "" {
		return errors.New("missing required input: QualificationCriteria")
	}
	if input.DisqualificationCriteria == "" {
		return errors.New("missing required input: DisqualificationCriteria")
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
	return &ICPQualificationResult{}
}

func NewICPQualificationCapability(postgres *postgres_repository.Repositories, aiService interfaces.AIService) *ICPQualificationCapability {
	return &ICPQualificationCapability{
		postgresRepositories: postgres,
		aiService:            aiService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapabilityExecution[ICPQualificationInput, ICPQualificationResult, NoConfig] = (*ICPQualificationCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                            = (*ICPQualificationCapability)(nil)
)

func (c *ICPQualificationCapability) Execute(ctx context.Context, data ICPQualificationInput, config NoConfig) (ICPQualificationResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.executeICPQualification")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	if err := c.ValidateInput(data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return ICPQualificationResult{}, err
	}
	if err := c.ValidateConfig(config); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return ICPQualificationResult{}, err
	}

	// get ICP definition

	// get all company context
	company, err := c.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, data.PrimaryDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return ICPQualificationResult{}, err
	}
	if company == nil {
		err := errors.New("Company not set in global orgs, cannot run ICP qualification")
		tracing.TraceErr(span, err)
		return ICPQualificationResult{}, err
	}

	// build prompt
	systemPrompt := `You are a world class company analyst.  Your objective is to determine whether the company I provide you fits our ideal customer profile or not.  I will provide you with three datasets: 1. a description of our ideal customer, 2. criteria that automatically disqualifies companies, and 3. all the context I have about the company, including their name, location, and several descriptions taken from their website and linkedin pages.  I need you to respond with two things: 1. true/false indicator on whether the company is an ICP fit or not, and 2. a short (less than 100 characters) description of your decision rationale.  Please output this as a pipe delimited string like this:  ICP_Fit: false | Rationale: this is my rationale.`

	content := fmt.Sprintf(`
    Company Name: %s
    Company Domain: %s
    Year Founded: %s
    Employee Count: %s
    Location: %s, %s, %s
    Industry NAICS Code: %s
    Industry Name: %s
    Company Description 1: %s
    Company Description 2: %s
    Company Description 3: %s
    Company Description 4: %s
    Company Description 5: %s
    `, company.Name, company.PrimaryDomain, company.YearFounded, company.EmployeeCount,
		company.City, company.Region, company.CountryA2, company.IndustryNaicsCode,
		company.IndustryNaicsName, company.Description, company.SourceDescription1, company.SourceDescription2, company.SourceDescription3, company.SourceDescription4)

	// askAI
	answer, err := c.aiService.AskAI(ctx, enum.AIModelAnthropicHaiku, systemPrompt, content)
	if err != nil {
		tracing.TraceErr(span, err)
		return ICPQualificationResult{}, err
	}

	// parse answer
	err = c.validateAnswer(ctx, answer)
	if err != nil {
		tracing.TraceErr(span, err)
		return ICPQualificationResult{}, err
	}

	icpFit, rationale := c.parseAnswer(ctx, *answer)

	// save to timeline

	// disqualify as lead if not a fit

	return ICPQualificationResult{}, nil
}

func (c *ICPQualificationCapability) validateAnswer(ctx context.Context, answer *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ICPQualificationCapability.parseAnswer")
	defer span.Finish()
	tracing.TagComponentService(span)
	span.LogKV("answer", utils.IfNotNilString(answer))

	if answer == nil {
		return errors.New("answer is empty")
	}

	if !strings.Contains(*answer, "ICP_Fit:") {
		return fmt.Errorf("missing ICP_Fit in input string")
	}
	if !strings.Contains(*answer, "| Rationale:") {
		return fmt.Errorf("missing Rationale in input string")
	}
	return nil
}

func (c *ICPQualificationCapability) parseAnswer(ctx context.Context, answer string) (bool, string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ICPQualificationCapability.parseAnswer")
	defer span.Finish()
	tracing.TagComponentService(span)

	parts := strings.Split(answer, " | Rationale: ")
	if len(parts) != 2 {
		return false, ""
	}

	// Get ICP_Fit value by trimming prefix and converting to bool
	icpFitStr := strings.TrimPrefix(parts[0], "ICP_Fit: ")
	icpFit := strings.ToLower(icpFitStr) == "true"

	// Get rationale
	rationale := parts[1]

	return icpFit, rationale
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
