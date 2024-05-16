package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
	"time"
)

type OrganizationService interface {
	WebScrapeOrganizations()
	RefreshLastTouchpoint()
	UpkeepOrganizations()
}

type organizationService struct {
	cfg                    *config.Config
	log                    logger.Logger
	repositories           *repository.Repositories
	eventsProcessingClient *grpc_client.Clients
}

func NewOrganizationService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, client *grpc_client.Clients) OrganizationService {
	return &organizationService{
		cfg:                    cfg,
		log:                    log,
		repositories:           repositories,
		eventsProcessingClient: client,
	}
}

func (s *organizationService) WebScrapeOrganizations() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	if s.eventsProcessingClient == nil {
		s.log.Warn("eventsProcessingClient is nil. Will not update next cycle date.")
		return
	}

	s.webScrapeOrganizations(ctx)
}

func (s *organizationService) webScrapeOrganizations(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.webScrapeOrganizations")
	defer span.Finish()

	records, err := s.repositories.OrganizationRepository.GetOrganizationsForWebScrape(ctx, s.cfg.ProcessConfig.WebScrapedOrganizationsPerCycle)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting organizations for status update: %v", err)
		return
	}
	span.LogFields(log.Int("organizations", len(records)))

	// no organizations found for web scraping
	if len(records) == 0 {
		return
	}

	// web scrape organizations
	for _, record := range records {
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return s.eventsProcessingClient.OrganizationClient.WebScrapeOrganization(ctx, &organizationpb.WebScrapeOrganizationGrpcRequest{
				Tenant:         record.Tenant,
				OrganizationId: record.OrganizationId,
				AppSource:      constants.AppSourceDataUpkeeper,
				Url:            record.Url,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error web scraping organization {%s}: %s", record.OrganizationId, err.Error())
		}
	}
}

func (s *organizationService) RefreshLastTouchpoint() {
	if s.eventsProcessingClient == nil {
		s.log.Warn("eventsProcessingClient is nil. Will not update next cycle date.")
		return
	}

	headers := map[string]string{
		"X-Openline-TENANT":  "openlineai",
		"X-Openline-API-KEY": s.cfg.PlatformAdminApi.ApiKey,
	}

	req, err := http.NewRequest("POST", s.cfg.PlatformAdminApi.Url+"/organization/refreshLastTouchpoint", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("RefreshLastTouchpoint: Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("RefreshLastTouchpoint: Error response:", resp.Status)
		return
	}
}

func (s *organizationService) UpkeepOrganizations() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	if s.eventsProcessingClient == nil {
		s.log.Warn("eventsProcessingClient is nil. Will not update next cycle date.")
		return
	}

	now := utils.Now()

	s.updateDerivedNextRenewalDates(ctx, now)
}

func (s *organizationService) updateDerivedNextRenewalDates(ctx context.Context, referenceTime time.Time) {
	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.updateDerivedNextRenewalDates")
	defer span.Finish()

	limit := 1000

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsForUpdateNextRenewalDate(ctx, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting contracts for renewal rollout: %v", err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		//process contracts
		for _, record := range records {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
				return s.eventsProcessingClient.OrganizationClient.RefreshRenewalSummary(ctx, &organizationpb.RefreshRenewalSummaryGrpcRequest{
					Tenant:         record.Tenant,
					OrganizationId: record.OrganizationId,
					AppSource:      constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error refreshing renewal summary for organization {%s}: %s", record.OrganizationId, err.Error())
			}
		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		//sleep for async processing, then check again
		time.Sleep(5 * time.Second)
	}
}
