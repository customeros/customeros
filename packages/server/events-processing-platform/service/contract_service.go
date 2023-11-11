package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contract"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go/log"
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
	span.LogFields(log.Object("request", request))

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

	source := commonmodel.Source{}
	source.FromGrpc(request.SourceFields)

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	// Determine the status based on ServiceStartedAt
	var contractStatus model.ContractStatus
	if request.ServiceStartedAt == nil || request.ServiceStartedAt.AsTime().After(utils.Now()) {
		contractStatus = model.Draft
	} else {
		contractStatus = model.Live
	}

	createContractCommand := command.NewCreateContractCommand(
		contractId,
		request.Tenant,
		request.LoggedInUserId,
		model.ContractDataFields{
			OrganizationId:   request.OrganizationId,
			Name:             request.Name,
			CreatedByUserId:  utils.StringFirstNonEmpty(request.CreatedByUserId, request.LoggedInUserId),
			ServiceStartedAt: utils.ToDatePtr(utils.TimestampProtoToTimePtr(request.ServiceStartedAt)),
			SignedAt:         utils.TimestampProtoToTimePtr(request.SignedAt),
			RenewalCycle:     model.RenewalCycle(request.RenewalCycle),
			Status:           contractStatus,
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
