package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/events_processing_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	orggrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"time"
)

type OrganizationService interface {
	UpdateNextCycleDate()
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

func (s *organizationService) UpdateNextCycleDate() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	if s.eventsProcessingClient == nil {
		s.log.Warn("eventsProcessingClient is nil. Will not update next cycle date.")
		return
	}

	span, ctx := tracing.StartTracerSpan(ctx, "OrganizationService.UpdateNextCycleDate")
	defer span.Finish()

	now := utils.Now()

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.OrganizationRepository.GetOrganizationsForNextCycleDateRenew(ctx, now)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting organizations for next cycle date renew: %v", err)
			return
		}

		// no organizations found for next cycle date renew
		if len(records) == 0 {
			return
		}

		//process organizations
		for _, record := range records {
			_, err = s.eventsProcessingClient.OrganizationClient.RequestRenewNextCycleDate(ctx, &orggrpc.RequestRenewNextCycleDateRequest{
				Tenant:         record.Tenant,
				OrganizationId: record.OrganizationId,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error requesting organization next cycle date renew: %s", err.Error())
			}
		}

		//sleep for async processing, then check again
		time.Sleep(5 * time.Second)
	}

}
