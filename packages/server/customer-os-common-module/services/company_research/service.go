package company_research

import (
	"context"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type companyResearchService struct {
	postgresRepositories *postgres_repository.Repositories
	aiService            interfaces.AIService
}

func NewCompanyResearchService(
	postgresRepositories *postgres_repository.Repositories,
	aiService interfaces.AIService,
) interfaces.CompanyResearch {
	return &companyResearchService{
		postgresRepositories: postgresRepositories,
		aiService:            aiService,
	}
}

func (s *companyResearchService) GenerateIdealCustomerProfile(ctx context.Context, tenantDomain string, trainingWebsites []string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "icpService.GenerateIdealCustomerProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	return "", nil
}

func (s *companyResearchService) GenerateCompanyBrief(ctx context.Context, domains []string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companyResearchService.GenerateCompanyBrief")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	return "", nil
}
