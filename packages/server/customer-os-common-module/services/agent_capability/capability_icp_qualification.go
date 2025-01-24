package agent_capability

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
)

type ICPQualificationCapability struct {
	postgresRepositories *postgres_repository.Repositories
	aiService            interfaces.AIService
}

func NewICPQualificationCapability(postgres *postgres_repository.Repositories, aiService interfaces.AIService) *ICPQualificationCapability {
	return &ICPQualificationCapability{
		postgresRepositories: postgres,
		aiService:            aiService,
	}
}

// Compile-time interface check
var _ interfaces.AgentCapabilityExecution[ICPQualificationInput, ICPQualificationResult] = (*ICPQualificationCapability)(nil)

type ICPQualificationInput struct {
	ICPDefinition string
	PrimaryDomain string
}

type ICPQualificationResult struct {
}

func (c *ICPQualificationCapability) Execute(ctx context.Context, data ICPQualificationInput) (ICPQualificationResult, error) {
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
