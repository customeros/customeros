package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/events_processing_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go/log"
)

type OrganizationService interface {
	WebScrapeOrganizations()
}

type organizationService struct {
	cfg                    *config.Config
	log                    logger.Logger
	repositories           *repository.Repositories
	eventsProcessingClient *events_processing_client.Client
}

func NewOrganizationService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, client *events_processing_client.Client) OrganizationService {
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
		_, err = s.eventsProcessingClient.OrganizationClient.WebScrapeOrganization(ctx, &organizationpb.WebScrapeOrganizationGrpcRequest{
			Tenant:         record.Tenant,
			OrganizationId: record.OrganizationId,
			AppSource:      constants.AppSourceDataUpkeeper,
			Url:            record.Url,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error web scraping organization {%s}: %s", record.OrganizationId, err.Error())
		}
	}
}
