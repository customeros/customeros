package service

import (
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/model"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type OrganizationService interface {
	DedupOrganizations()
}

type organizationService struct {
	cfg           *config.Config
	log           logger.Logger
	repositories  *repository.Repositories
	graphqlClient *graphql.Client
}

func NewOrganizationService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, graphqlClient *graphql.Client) OrganizationService {
	return &organizationService{
		cfg:           cfg,
		log:           log,
		repositories:  repositories,
		graphqlClient: graphqlClient,
	}
}

type OrganizationIdName struct {
	Id   string
	Name string
}

func (s *organizationService) DedupOrganizations() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel() // Cancel context on exit

	tenants, err := s.getTenantsWithOrganizations(ctx)
	if err != nil {
		s.log.Errorf("Failed to get tenants: %v", err)
	} else {
		s.log.Infof("Got %d tenants for organization dedup", len(tenants))
	}

	// Long-running dedup
	for _, tenant := range tenants {
		s.dedupTenantOrganizations(ctx, tenant)
	}
}

func (s *organizationService) dedupTenantOrganizations(ctx context.Context, tenant string) {
	span, ctx := tracing.StartTracerSpan(ctx, "dedupTenantOrganizations")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))

	pageSize := s.cfg.Organizations.ApiPageSize
	windowSize := s.cfg.Organizations.CompareWindowSize
	window := make([]OrganizationIdName, 0, windowSize)

	page := 1
	for {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return
		default:
			// Continue fetching organizations
		}

		// Call API to get organizations
		orgs, hasMore, err := s.getOrganizationsPage(ctx, tenant, page, pageSize)
		if err != nil {
			s.log.Errorf("Failed to get organizations: %v", err)
			tracing.TraceErr(span, err)
			return
		}

		// Append page to window
		window = append(window, orgs...)

		// Keep window sized
		if len(window) > windowSize {
			excess := len(window) - windowSize
			window = window[excess:]
		}

		// Compare organizations
		//compareOrgs(window)

		// Break if no more pages
		if !hasMore {
			break
		}

		// Next page
		page++
	}
}

func (s *organizationService) getOrganizationsPage(ctx context.Context, tenant string, page, pageSize int) ([]OrganizationIdName, bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.getOrganizationsPage")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant), log.Int("page", page), log.Int("pageSize", pageSize))

	graphqlRequest := graphql.NewRequest(
		`query Organizations($page: Int!, $pageSize: Int!) {
  				organizations(pagination: {limit: $pageSize, page: $page},
					sort: {by: "NAME" caseSensitive: false}) {
				content {
      				id
      				name
    			}}
			}`)
	graphqlRequest.Var("page", page)
	graphqlRequest.Var("pageSize", pageSize)
	s.addHeadersToGraphRequest(graphqlRequest, tenant)

	var graphqlResponse model.OrganizationsResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed query cosApi :%v", err.Error())
		return nil, false, err
	}

	orgs := make([]OrganizationIdName, len(graphqlResponse.Organizations.Content))
	for i, org := range graphqlResponse.Organizations.Content {
		orgs[i] = OrganizationIdName{
			Id:   org.ID,
			Name: org.Name,
		}
	}

	hasMore := len(graphqlResponse.Organizations.Content) == pageSize

	return orgs, hasMore, nil
}

func (s *organizationService) getTenantsWithOrganizations(ctx context.Context) ([]string, error) {
	return s.repositories.TenantRepository.GetTenantsWithOrganizations(ctx, s.cfg.Organizations.AtLeastPerTenant)
}

func (s *organizationService) addHeadersToGraphRequest(req *graphql.Request, tenant string) {
	req.Header.Add("X-Openline-API-KEY", s.cfg.Service.CustomerOsAdminAPIKey)
	if tenant != "" {
		req.Header.Add("X-Openline-TENANT", tenant)
	}
}
