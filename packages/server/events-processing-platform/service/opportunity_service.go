package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/opportunity"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	organizationaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
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

type opportunityService struct {
	opportunitypb.UnimplementedOpportunityGrpcServiceServer
	log                        logger.Logger
	opportunityCommandHandlers *command_handler.CommandHandlers
	aggregateStore             eventstore.AggregateStore
}

func NewOpportunityService(log logger.Logger, commandHandlers *command_handler.CommandHandlers, aggregateStore eventstore.AggregateStore) *opportunityService {
	return &opportunityService{
		log:                        log,
		opportunityCommandHandlers: commandHandlers,
		aggregateStore:             aggregateStore,
	}
}

func (s *opportunityService) CreateOpportunity(ctx context.Context, request *opportunitypb.CreateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OpportunityService.CreateOpportunity")
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

	opportunityId := uuid.New().String()

	// Convert any protobuf timestamp to time.Time, if necessary
	createdAt := utils.TimestampProtoToTimePtr(request.CreatedAt)
	updatedAt := utils.TimestampProtoToTimePtr(request.UpdatedAt)
	estimatedClosedAt := utils.TimestampProtoToTimePtr(request.EstimatedCloseDate)

	source := commonmodel.Source{}
	source.FromGrpc(request.SourceFields)

	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	createOpportunityCommand := command.NewCreateOpportunityCommand(
		opportunityId,
		request.Tenant,
		request.LoggedInUserId,
		model.OpportunityDataFields{
			Name:              request.Name,
			Amount:            request.Amount,
			InternalType:      model.OpportunityInternalType(request.InternalType),
			ExternalType:      request.ExternalType,
			InternalStage:     model.OpportunityInternalStage(request.InternalStage),
			ExternalStage:     request.ExternalStage,
			EstimatedClosedAt: estimatedClosedAt,
			OwnerUserId:       request.OwnerUserId,
			CreatedByUserId:   utils.StringFirstNonEmpty(request.CreatedByUserId, request.LoggedInUserId),
			GeneralNotes:      request.GeneralNotes,
			NextSteps:         request.NextSteps,
			OrganizationId:    request.OrganizationId,
		},
		source,
		externalSystem,
		createdAt,
		updatedAt,
	)

	if err := s.opportunityCommandHandlers.CreateOpportunity.Handle(ctx, createOpportunityCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateOpportunity.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created opportunity
	return &opportunitypb.OpportunityIdGrpcResponse{Id: opportunityId}, nil
}

func (s *opportunityService) UpdateRenewalOpportunity(ctx context.Context, request *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OpportunityService.UpdateRenewalOpportunity")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	span.LogFields(log.Object("request", request))

	// Check if the opportunity ID is valid
	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	// Convert any protobuf timestamp to time.Time, if necessary
	updatedAt := utils.TimestampProtoToTimePtr(request.UpdatedAt)

	source := commonmodel.Source{}
	source.FromGrpc(request.SourceFields)

	updateRenewalOpportunityCommand := command.NewUpdateRenewalOpportunityCommand(
		request.Id,
		request.Tenant,
		request.LoggedInUserId,
		request.Comments,
		model.RenewalLikelihood(request.RenewalLikelihood).StringValue(),
		request.Amount,
		source,
		updatedAt,
	)

	if err := s.opportunityCommandHandlers.UpdateRenewalOpportunity.Handle(ctx, updateRenewalOpportunityCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateRenewalOpportunity.Handle) tenant:{%v}, opportunityId:{%v}, err: %v", request.Tenant, request.Id, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created opportunity
	return &opportunitypb.OpportunityIdGrpcResponse{Id: request.Id}, nil
}

func (s *opportunityService) checkOrganizationExists(ctx context.Context, tenant, organizationId string) (bool, error) {
	organizationAggregate := organizationaggregate.NewOrganizationAggregateWithTenantAndID(tenant, organizationId)
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
