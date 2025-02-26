package company_research

import (
	"context"
	"errors"
	"fmt"
	"strings"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type companyResearchService struct {
	postgresRepositories *postgres_repository.Repositories
	aiService            interfaces.AIService
	webscraperService    interfaces.WebscraperService
	workspaceService     interfaces.WorkspaceService
}

func NewCompanyResearchService(
	postgresRepositories *postgres_repository.Repositories,
	aiService interfaces.AIService,
	webscraperService interfaces.WebscraperService,
	workspaceService interfaces.WorkspaceService,
) interfaces.CompanyResearch {
	return &companyResearchService{
		postgresRepositories: postgresRepositories,
		aiService:            aiService,
		webscraperService:    webscraperService,
		workspaceService:     workspaceService,
	}
}

const (
	WebpagesInChunk = 25
	AIModel         = enum.AIModelAnthropicSonnet
)

func (s *companyResearchService) GenerateIdealCustomerProfile(ctx context.Context, tenantDomain string, trainingWebsites []string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "icpService.GenerateIdealCustomerProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	return "", nil
}

func (s *companyResearchService) GenerateCompanyBriefForTenant(ctx context.Context) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companyResearchService.GenerateCompanyBriefForTenant")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	domains, err := s.workspaceService.GetWorkspaceDomainsForTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	report, err := s.GenerateCompanyBrief(ctx, domains)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// save report to db
	err = s.postgresRepositories.TenantRepository.SetCompanyReport(ctx, *report)
	if err != nil {
		tracing.TraceErr(span, err)
		return report, err
	}

	return report, nil
}

func (s *companyResearchService) GenerateCompanyBrief(ctx context.Context, domains []string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companyResearchService.GenerateCompanyBrief")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// crawl domains
	for _, domain := range domains {
		url := "https://" + domain
		_, err := s.webscraperService.Crawl(ctx, url)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	// get all webpages
	webpages, err := s.postgresRepositories.GlobalOrganizationWebpageRepository.GetAllWebpagesByPrimaryDomains(ctx, domains)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// generate company report
	report, err := s.generateCompanyBrief(ctx, domains, webpages)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if report == nil {
		return nil, errors.New("unable to generate company brief")
	}

	return report, nil
}

func (s *companyResearchService) generateCompanyBrief(ctx context.Context, domains []string, webpages []*postgres_entity.GlobalOrganizationWebpages) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companyResearchService.generateCompanyBrief")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	chunks := s.buildDomainChunks(webpages)

	var intermediateReports []string
	for _, domain := range domains {
		globalOrgDetails, err := s.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, err)
		}

		for i, chunk := range chunks[domain] {
			systemPrompt, prompt := s.buildChunkPrompt(domain, chunk, i+1, *globalOrgDetails)
			answer, err := s.aiService.AskAI(ctx, AIModel, systemPrompt, prompt)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
			intermediateReports = append(intermediateReports, *answer)
		}
	}

	report, err := s.buildFinalReport(ctx, intermediateReports)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if report == nil {
		return nil, errors.New("unable to produce company brief")
	}

	return report, nil
}

func (s *companyResearchService) buildFinalReport(ctx context.Context, analyses []string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "companyResearchService.buildFinalReport")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Construct synthesis prompt
	var sp strings.Builder
	sp.WriteString("Based on the following analyses of scraped web content, create a comprehensive company report. ")
	sp.WriteString("The analysis may come from different domains.  If that's the case, all domains belong to the same company.")
	sp.WriteString("Structure the report with these sections:\n")
	sp.WriteString("1. Executive Overview\n")
	sp.WriteString("2. Market Position\n")
	sp.WriteString("3. Product Analysis\n")
	sp.WriteString("4. Customer Profile\n")
	sp.WriteString("5. Growth Indicators\n")
	sp.WriteString("6. Risk Assessment\n\n")
	sp.WriteString("Create a professional, comprehensive report with specific details and insights.")
	sp.WriteString("Please respond in valid markdown format with only the actual report content.  No additional text or descriptions.")

	var p strings.Builder
	for i, analysis := range analyses {
		p.WriteString(fmt.Sprintf("--- ANALYSIS %d ---\n", i+1))
		p.WriteString(analysis)
		p.WriteString("\n\n")
	}

	report, err := s.aiService.AskAI(ctx, AIModel, sp.String(), p.String())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return report, nil
}

func (s *companyResearchService) buildChunkPrompt(domain string, webpageChunk []*postgres_entity.GlobalOrganizationWebpages, chunkCount int, globalOrgDetails postgres_entity.GlobalOrganization) (string, string) {
	var systemPrompt strings.Builder

	systemPrompt.WriteString("You are analyzing scraped web content to create a detailed company report. ")
	systemPrompt.WriteString(fmt.Sprintf("This is chunk %d of the analysis. ", chunkCount))
	systemPrompt.WriteString("Extract key information about the company's: ")
	systemPrompt.WriteString("1. Products/services\n")
	systemPrompt.WriteString("2. Target markets\n")
	systemPrompt.WriteString("3. Value propositions\n")
	systemPrompt.WriteString("4. Competitive advantages\n")
	systemPrompt.WriteString("5. Customer types\n")
	systemPrompt.WriteString("6. Growth indicators\n\n")
	systemPrompt.WriteString("Please respond in valid markdown format with only the actual report content.  No additional text or descriptions.")
	systemPrompt.WriteString("Web page content follows:\n\n")

	var prompt strings.Builder

	prompt.WriteString("--- BASIC COMPANY FACTS ---\n")
	prompt.WriteString(fmt.Sprintf("Company Name: %s\n", globalOrgDetails.Name))
	prompt.WriteString(fmt.Sprintf("Year Founded: %d\n", globalOrgDetails.YearFounded))
	prompt.WriteString(fmt.Sprintf("Employee Count: %d\n", globalOrgDetails.EmployeeCount))
	prompt.WriteString(fmt.Sprintf("Location: %s, %s, %s\n", globalOrgDetails.City, globalOrgDetails.Region, globalOrgDetails.CountryA2))
	prompt.WriteString(fmt.Sprintf("Industry Name: %s\n", globalOrgDetails.IndustryNaicsName))
	prompt.WriteString("\n\n")

	for i, page := range webpageChunk {
		sections := strings.Split(page.Content, "Links/Buttons")
		prompt.WriteString(fmt.Sprintf("--- PAGE %d: ---\n", i+1))
		prompt.WriteString(fmt.Sprintf("URL: %s\n", page.Url))
		prompt.WriteString(sections[0])
		prompt.WriteString("\n\n")
	}

	prompt.WriteString("\nProvide a structured analysis with specific facts and quotes from the content.")

	return systemPrompt.String(), prompt.String()
}

func (s *companyResearchService) buildDomainChunks(webpages []*postgres_entity.GlobalOrganizationWebpages) map[string][][]*postgres_entity.GlobalOrganizationWebpages {
	// Group webpages by domain
	domainMap := make(map[string][]*postgres_entity.GlobalOrganizationWebpages)

	for _, webpage := range webpages {
		domain := webpage.PrimaryDomain
		domainMap[domain] = append(domainMap[domain], webpage)
	}

	// Create chunks
	result := make(map[string][][]*postgres_entity.GlobalOrganizationWebpages)

	for domain, pages := range domainMap {
		var chunks [][]*postgres_entity.GlobalOrganizationWebpages

		// Split into chunks
		for i := 0; i < len(pages); i += WebpagesInChunk {
			end := i + WebpagesInChunk
			if end > len(pages) {
				end = len(pages)
			}

			chunks = append(chunks, pages[i:end])
		}

		result[domain] = chunks
	}

	return result
}
