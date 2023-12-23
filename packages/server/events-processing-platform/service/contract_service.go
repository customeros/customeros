package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contractService struct {
	contractpb.UnimplementedContractGrpcServiceServer
	log                     logger.Logger
	contractCommandHandlers *command_handler.CommandHandlers
	aggregateStore          eventstore.AggregateStore
}

func NewContractService(log logger.Logger, commandHandlers *command_handler.CommandHandlers, aggregateStore eventstore.AggregateStore) *contractService {
	return &contractService{
		log:                     log,
		contractCommandHandlers: commandHandlers,
		aggregateStore:          aggregateStore,
	}
}

func (s *contractService) CreateContract(ctx context.Context, request *contractpb.CreateContractGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContractService.CreateContract")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	// Validate organization ID
	if request.OrganizationId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("organizationId"))
	}
	// Check if the organization aggregate exists
	orgExists, err := s.checkOrganizationExists(ctx, request.Tenant, request.OrganizationId)
	if err != nil {
		s.log.Error(err, "error checking organization existence")
		return nil, status.Errorf(codes.Internal, "error checking organization existence: %v", err)
	}
	if !orgExists {
		return nil, status.Errorf(codes.NotFound, "organization with ID %s not found", request.OrganizationId)
	}

	contractId := uuid.New().String()

	// Convert any protobuf timestamp to time.Time, if necessary
	createdAt, updatedAt := convertCreateAndUpdateProtoTimestampsToTime(request.CreatedAt, request.UpdatedAt)
	serviceStartedAt := utils.TimestampProtoToTimePtr(request.ServiceStartedAt)
	signedAt := utils.TimestampProtoToTimePtr(request.SignedAt)

	source := commonmodel.Source{}
	source.FromGrpc(request.SourceFields)

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	createContractCommand := command.NewCreateContractCommand(
		contractId,
		request.Tenant,
		request.LoggedInUserId,
		model.ContractDataFields{
			OrganizationId:   request.OrganizationId,
			Name:             request.Name,
			ContractUrl:      request.ContractUrl,
			CreatedByUserId:  utils.StringFirstNonEmpty(request.CreatedByUserId, request.LoggedInUserId),
			ServiceStartedAt: serviceStartedAt,
			SignedAt:         signedAt,
			RenewalCycle:     model.RenewalCycle(request.RenewalCycle),
			RenewalPeriods:   request.RenewalPeriods,
		},
		source,
		externalSystem,
		createdAt,
		updatedAt,
	)

	if err := s.contractCommandHandlers.CreateContract.Handle(ctx, createContractCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateContract.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created contract
	return &contractpb.ContractIdGrpcResponse{Id: contractId}, nil
}

func (s *contractService) UpdateContract(ctx context.Context, request *contractpb.UpdateContractGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContractService.UpdateContract")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	// Check if the contract ID is valid
	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	// Convert any protobuf timestamp to time.Time, if necessary
	updatedAt := utils.TimestampProtoToTimePtr(request.UpdatedAt)
	serviceStartedAt := utils.TimestampProtoToTimePtr(request.ServiceStartedAt)
	signedAt := utils.TimestampProtoToTimePtr(request.SignedAt)
	endedAt := utils.TimestampProtoToTimePtr(request.EndedAt)

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	source := commonmodel.Source{}
	source.FromGrpc(request.SourceFields)

	// Create update contract command
	updateContractCommand := command.NewUpdateContractCommand(
		request.Id,
		request.Tenant,
		request.LoggedInUserId,
		model.ContractDataFields{
			Name:             request.Name,
			ServiceStartedAt: serviceStartedAt,
			SignedAt:         signedAt,
			EndedAt:          endedAt,
			RenewalCycle:     model.RenewalCycle(request.RenewalCycle),
			ContractUrl:      request.ContractUrl,
			RenewalPeriods:   request.RenewalPeriods,
		},
		source,
		externalSystem,
		updatedAt,
	)

	if err := s.contractCommandHandlers.UpdateContract.Handle(ctx, updateContractCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateContract.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &contractpb.ContractIdGrpcResponse{Id: request.Id}, nil
}

func (s *contractService) RefreshContractStatus(ctx context.Context, request *contractpb.RefreshContractStatusGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContractService.RefreshContractStatus")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	cmd := command.NewRefreshContractStatusCommand(request.Id, request.Tenant, request.LoggedInUserId, request.AppSource)

	if err := s.contractCommandHandlers.RefreshContractStatus.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RefreshContractStatus.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &contractpb.ContractIdGrpcResponse{Id: request.Id}, nil
}

func (s *contractService) RolloutRenewalOpportunityOnExpiration(ctx context.Context, request *contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest) (*contractpb.ContractIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ContractService.RefreshContractStatus")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	cmd := command.NewRolloutRenewalOpportunityOnExpirationCommand(request.Id, request.Tenant, request.LoggedInUserId, request.AppSource)

	if err := s.contractCommandHandlers.RolloutRenewalOpportunityOnExpiration.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(RolloutRenewalOpportunityOnExpiration.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &contractpb.ContractIdGrpcResponse{Id: request.Id}, nil
}

func (s *contractService) checkOrganizationExists(ctx context.Context, tenant, organizationId string) (bool, error) {
	organizationAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenant, organizationId)
	err := s.aggregateStore.Exists(ctx, organizationAggregate.GetID())
	if err != nil {
		if errors.Is(err, eventstore.ErrAggregateNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil // The organization exists
}
