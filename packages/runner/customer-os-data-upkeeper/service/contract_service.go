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
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contract"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/opportunity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// this is a catch-all for contracts that have ended but still have active renewal opportunities
	s.closeEndedContractOpportunityRenewals(ctx, now)
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

		// no contracts found for next cycle date renew
		if len(records) == 0 {
			return
		}

		//process contracts
		for _, record := range records {
			_, err = s.eventsProcessingClient.ContractClient.RefreshContractStatus(ctx, &contractpb.RefreshContractStatusGrpcRequest{
				Tenant:    record.Tenant,
				Id:        record.ContractId,
				AppSource: constants.AppSourceDataUpkeeper,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error refreshing contract status: %s", err.Error())
				grpcErr, ok := status.FromError(err)
				if ok && grpcErr.Code() == codes.NotFound && grpcErr.Message() == "aggregate not found" {
					s.resyncContract(ctx, record.Tenant, record.ContractId)
				}
			}
		}

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

		// no contracts found for next cycle date renew
		if len(records) == 0 {
			return
		}

		//process contracts
		for _, record := range records {
			_, err = s.eventsProcessingClient.ContractClient.RolloutRenewalOpportunityOnExpiration(ctx, &contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest{
				Tenant:    record.Tenant,
				Id:        record.ContractId,
				AppSource: constants.AppSourceDataUpkeeper,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error rollout renewal opportunity: %s", err.Error())
				grpcErr, ok := status.FromError(err)
				if ok && grpcErr.Code() == codes.NotFound && grpcErr.Message() == "aggregate not found" {
					s.resyncContract(ctx, record.Tenant, record.ContractId)
				}
			}
		}

		//sleep for async processing, then check again
		time.Sleep(5 * time.Second)
	}
}

func (s *contractService) closeEndedContractOpportunityRenewals(ctx context.Context, referenceTime time.Time) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContractService.closeEndedContractOpportunityRenewals")
	defer span.Finish()

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.OpportunityRepository.GetRenewalOpportunitiesForClosingAsLost(ctx, referenceTime)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting opportunities for closing: %v", err)
			return
		}

		// no renewal opportunities found, return
		if len(records) == 0 {
			return
		}

		//process renewal opportunities
		for _, record := range records {
			_, err = s.eventsProcessingClient.OpportunityClient.CloseLooseOpportunity(ctx, &opportunitypb.CloseLooseOpportunityGrpcRequest{
				Tenant:    record.Tenant,
				Id:        record.OpportunityId,
				AppSource: constants.AppSourceDataUpkeeper,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error closing renewal opportunity: %s", err.Error())
			}
		}

		//sleep for async processing, then check again
		time.Sleep(5 * time.Second)
	}
}

func (s *contractService) resyncContract(ctx context.Context, tenant, contractId string) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContractService.resyncContract")
	defer span.Finish()

	contractDbNode, err := s.repositories.ContractRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting contract {%s}: %s", contractId, err.Error())
		return
	}

	props := utils.GetPropsFromNode(*contractDbNode)

	request := contractpb.UpdateContractGrpcRequest{
		Tenant:           tenant,
		Id:               contractId,
		Name:             utils.GetStringPropOrEmpty(props, "name"),
		ContractUrl:      utils.GetStringPropOrEmpty(props, "contractUrl"),
		ServiceStartedAt: utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(props, "serviceStartedAt")),
		SignedAt:         utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(props, "signedAt")),
		EndedAt:          utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(props, "endedAt")),
		SourceFields: &commonpb.SourceFields{
			Source:    utils.GetStringPropOrEmpty(props, "sourceOfTruth"),
			AppSource: constants.AppSourceDataUpkeeper,
		},
	}
	switch utils.GetStringPropOrEmpty(props, "renewalCycle") {
	case "MONTHLY":
		request.RenewalCycle = contractpb.RenewalCycle_MONTHLY_RENEWAL
	case "ANNUALLY":
		request.RenewalCycle = contractpb.RenewalCycle_ANNUALLY_RENEWAL
	default:
		request.RenewalCycle = contractpb.RenewalCycle_NONE
	}
	_, err = s.eventsProcessingClient.ContractClient.UpdateContract(ctx, &request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error re-syncing contract {%s}: %s", contractId, err.Error())
	}
}
