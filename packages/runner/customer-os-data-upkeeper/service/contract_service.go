package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type ContractService interface {
	UpkeepContracts()
	resyncContract(ctx context.Context, tenant, contractId string)
}

type contractService struct {
	cfg                    *config.Config
	log                    logger.Logger
	repositories           *repository.Repositories
	eventsProcessingClient *grpc_client.Clients
	services               *service.Services
}

func NewContractService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, client *grpc_client.Clients, services *service.Services) ContractService {
	return &contractService{
		cfg:                    cfg,
		log:                    log,
		repositories:           repositories,
		eventsProcessingClient: client,
		services:               services,
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
	s.closeActiveRenewalOpportunitiesForEndedContracts(ctx)
	s.createRenewalOpportunitiesIfMissing(ctx)
}

func (s *contractService) updateContractStatuses(ctx context.Context, referenceTime time.Time) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContractService.updateContractStatuses")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 100

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractsForStatusRenewal(ctx, referenceTime, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting contracts for status update: %v", err)
			return
		}

		// no contracts found for status update
		if len(records) == 0 {
			return
		}

		//process contracts
		for _, record := range records {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*contractpb.ContractIdGrpcResponse](func() (*contractpb.ContractIdGrpcResponse, error) {
				return s.eventsProcessingClient.ContractClient.RefreshContractStatus(ctx, &contractpb.RefreshContractStatusGrpcRequest{
					Tenant:    record.Tenant,
					Id:        record.ContractId,
					AppSource: constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error refreshing contract status: %s", err.Error())
				grpcErr, ok := status.FromError(err)
				if ok && grpcErr.Code() == codes.NotFound && grpcErr.Message() == "aggregate not found" {
					s.resyncContract(ctx, record.Tenant, record.ContractId)
				}
			} else {
				err = s.repositories.Neo4jRepositories.ContractWriteRepository.MarkStatusRenewalRequested(ctx, record.Tenant, record.ContractId)
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Error marking status renewal requested: %s", err.Error())
				}
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

func (s *contractService) rolloutContractRenewals(ctx context.Context, referenceTime time.Time) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContractService.rolloutContractRenewals")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 100

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractsForRenewalRollout(ctx, referenceTime, limit)
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
			_, err = utils.CallEventsPlatformGRPCWithRetry[*contractpb.ContractIdGrpcResponse](func() (*contractpb.ContractIdGrpcResponse, error) {
				return s.eventsProcessingClient.ContractClient.RolloutRenewalOpportunityOnExpiration(ctx, &contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest{
					Tenant:    record.Tenant,
					Id:        record.ContractId,
					AppSource: constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error rollout renewal opportunity: %s", err.Error())
				grpcErr, ok := status.FromError(err)
				if ok && grpcErr.Code() == codes.NotFound && grpcErr.Message() == "aggregate not found" {
					s.resyncContract(ctx, record.Tenant, record.ContractId)
				}
			} else {
				err = s.repositories.Neo4jRepositories.ContractWriteRepository.MarkRolloutRenewalRequested(ctx, record.Tenant, record.ContractId)
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Error marking renewal rollout requested: %s", err.Error())
				}
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

func (s *contractService) closeActiveRenewalOpportunitiesForEndedContracts(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContractService.closeActiveRenewalOpportunitiesForEndedContracts")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 100

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.OpportunityReadRepository.GetRenewalOpportunitiesForClosingAsLost(ctx, limit)
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
			_, err = utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
				return s.eventsProcessingClient.OpportunityClient.CloseLooseOpportunity(ctx, &opportunitypb.CloseLooseOpportunityGrpcRequest{
					Tenant:    record.Tenant,
					Id:        record.OpportunityId,
					AppSource: constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error closing renewal opportunity: %s", err.Error())
			} else {
				err = s.repositories.Neo4jRepositories.OpportunityWriteRepository.MarkRenewalRequested(ctx, record.Tenant, record.OpportunityId)
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Error marking renewal rollout requested: %s", err.Error())
				}
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

func (s *contractService) createRenewalOpportunitiesIfMissing(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContractService.createRenewalOpportunitiesIfMissing")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 100

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetLiveContractsWithoutRenewalOpportunities(ctx, limit)
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
			_, err = utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
				return s.eventsProcessingClient.OpportunityClient.CreateRenewalOpportunity(ctx, &opportunitypb.CreateRenewalOpportunityGrpcRequest{
					Tenant:     record.Tenant,
					ContractId: record.ContractId,
					SourceFields: &commonpb.SourceFields{
						AppSource: constants.AppSourceDataUpkeeper,
					},
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error creating renewal opportunity: %s", err.Error())
			}
		}

		// process only single batch per cycle
		return
	}
}

func (s *contractService) resyncContract(ctx context.Context, tenant, contractId string) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContractService.ResyncContract")
	defer span.Finish()

	contractDbNode, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting contract {%s}: %s", contractId, err.Error())
		return
	}

	props := utils.GetPropsFromNode(*contractDbNode)

	innerCtx := common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    tenant,
		AppSource: constants.AppSourceDataUpkeeper,
	})

	contractDataFields := data_fields.ContractSaveFields{
		Name:                 utils.StringPtr(utils.GetStringPropOrEmpty(props, "name")),
		ContractUrl:          utils.StringPtr(utils.GetStringPropOrEmpty(props, "contractUrl")),
		ServiceStartedAt:     utils.GetTimePropOrNil(props, "serviceStartedAt"),
		SignedAt:             utils.GetTimePropOrNil(props, "signedAt"),
		EndedAt:              utils.GetTimePropOrNil(props, "endedAt"),
		InvoicingStartDate:   utils.GetTimePropOrNil(props, "invoicingStartDate"),
		LengthInMonths:       utils.Int64Ptr(utils.GetInt64PropOrZero(props, "lengthInMonths")),
		BillingCycleInMonths: utils.Int64Ptr(utils.GetInt64PropOrZero(props, "billingCycleInMonths")),
	}
	currency := utils.GetStringPropOrEmpty(props, "currency")
	if currency != "" {
		contractDataFields.Currency = utils.ToPtr(neo4jenum.DecodeCurrency(currency))
	}

	_, err = s.services.ContractService.Save(innerCtx, &contractId, contractDataFields)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error re-syncing contract {%s}: %s", contractId, err.Error())
	}
}
