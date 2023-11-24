package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/events_processing_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contract"
	"time"
)

type ContractService interface {
	UpkeepContracts()
}

type contractService struct {
	cfg                    *config.Config
	log                    logger.Logger
	repositories           *repository.Repositories
	eventsProcessingClient *events_processing_client.Client
}

func NewContractService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, client *events_processing_client.Client) ContractService {
	return &contractService{
		cfg:                    cfg,
		log:                    log,
		repositories:           repositories,
		eventsProcessingClient: client,
	}
}

func (s *contractService) UpkeepContracts() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	if s.eventsProcessingClient == nil {
		s.log.Warn("eventsProcessingClient is nil. Will not update next cycle date.")
		return
	}

	now := utils.Now()

	s.updateContractStatuses(ctx, now)
	s.rolloutContractRenewals(ctx, now)
}

func (s *contractService) updateContractStatuses(ctx context.Context, referenceTime time.Time) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContractService.updateContractStatuses")
	defer span.Finish()

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.ContractRepository.GetContractsForStatusRenewal(ctx, referenceTime)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting contracts for status update: %v", err)
			return
		}

		// no organizations found for next cycle date renew
		if len(records) == 0 {
			return
		}

		//process organizations
		//for _, record := range records {
		//	_, err = s.eventsProcessingClient.ContractClient.RefreshContractStatus(ctx, &contractpb.RefreshContractStatusGrpcRequest{
		//		Tenant:    record.Tenant,
		//		Id:        record.ContractId,
		//		AppSource: constants.AppSourceDataUpkeeper,
		//	})
		//	if err != nil {
		//		tracing.TraceErr(span, err)
		//		s.log.Errorf("Error refreshing contract status: %s", err.Error())
		//	}
		//}

		//sleep for async processing, then check again
		time.Sleep(5 * time.Second)
	}
}

func (s *contractService) rolloutContractRenewals(ctx context.Context, referenceTime time.Time) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContractService.rolloutContractRenewals")
	defer span.Finish()

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.ContractRepository.GetContractsForRenewalRollout(ctx, referenceTime)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting contracts for renewal rollout: %v", err)
			return
		}

		// no organizations found for next cycle date renew
		if len(records) == 0 {
			return
		}

		//process organizations
		for _, record := range records {
			_, err = s.eventsProcessingClient.ContractClient.RolloutRenewalOpportunityOnExpiration(ctx, &contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest{
				Tenant:    record.Tenant,
				Id:        record.ContractId,
				AppSource: constants.AppSourceDataUpkeeper,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error refreshing contract status: %s", err.Error())
			}
		}

		//sleep for async processing, then check again
		time.Sleep(5 * time.Second)
	}
}
