package agent_capability

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
)

type ICPQualificationCapability struct {
	postgresRepositories *postgres_repository.Repositories
	aiService            interfaces.AIService
}

func (c *ICPQualificationCapability) ValidateConfig(config NoConfig) error {
	//TODO implement me
	panic("implement me")
}

func (c *ICPQualificationCapability) ValidateInput(input ICPQualificationInput) error {
	//TODO implement me
	panic("implement me")
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

type ICPQualificationInput struct {
	ICPDefinition string `json:"icpDefinition"`
	PrimaryDomain string `json:"primaryDomain"`
}

type ICPQualificationResult struct {
}

func (c *ICPQualificationCapability) Execute(ctx context.Context, data ICPQualificationInput, config NoConfig) (ICPQualificationResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.executeICPQualification")
	defer span.Finish()
	tracing.TagComponentService(span)

	// get ICP definition

	// get all company context
	_, err := c.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, data.PrimaryDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return ICPQualificationResult{}, err
	}

	// build prompt

	// askAI

	// save to timeline

	// disqualify as lead if not a fit

	return ICPQualificationResult{}, nil
}

// ExecuteUntyped implements the AgentCapabilityUntyped interface.
// It casts the generic input and config to the specific types and delegates to the typed Execute method.
func (c *ICPQualificationCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
	typedInput, ok := input.(*ICPQualificationInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type: expected ICPQualificationInput")
	}

	typedConfig, ok := config.(*NoConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type: expected NoCOnfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
