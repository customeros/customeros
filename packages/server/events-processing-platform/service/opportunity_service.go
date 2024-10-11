package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	contractaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	organizationaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type opportunityService struct {
	opportunitypb.UnimplementedOpportunityGrpcServiceServer
	services                   *Services
	log                        logger.Logger
	opportunityCommandHandlers *command_handler.CommandHandlers
	aggregateStore             eventstore.AggregateStore
}

func NewOpportunityService(log logger.Logger, commandHandlers *command_handler.CommandHandlers, aggregateStore eventstore.AggregateStore, services *Services) *opportunityService {
	return &opportunityService{
		services:                   services,
		log:                        log,
		opportunityCommandHandlers: commandHandlers,
		aggregateStore:             aggregateStore,
	}
}

func (s *opportunityService) CreateOpportunity(ctx context.Context, request *opportunitypb.CreateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OpportunityService.CreateOpportunity")
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

	opportunityId := uuid.New().String()
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewOpportunityAggregateWithTenantAndID(request.Tenant, opportunityId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateRenewalOpportunity) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created opportunity
	return &opportunitypb.OpportunityIdGrpcResponse{Id: opportunityId}, nil
}

func (s *opportunityService) UpdateOpportunity(ctx context.Context, request *opportunitypb.UpdateOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OpportunityService.UpdateOpportunity")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.Id)

	// Check if the opportunity ID is valid
	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewOpportunityAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateRenewalOpportunity) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &opportunitypb.OpportunityIdGrpcResponse{Id: request.Id}, nil
}

func (s *opportunityService) CreateRenewalOpportunity(ctx context.Context, request *opportunitypb.CreateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OpportunityService.CreateRenewalOpportunity")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	// Validate contract ID
	if request.ContractId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("contractId"))
	}

	// Check if the contract aggregate exists prior to closing the service line item
	contractExists, err := s.checkContractExists(ctx, request.Tenant, request.ContractId)
	if err != nil {
		s.log.Error(err, "error checking contract existence")
		return nil, status.Errorf(codes.Internal, "error checking contract existence: %v", err)
	}
	if !contractExists {
		return nil, status.Errorf(codes.NotFound, "contract with ID %s not found", request.ContractId)
	}

	opportunityId := uuid.New().String()
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewOpportunityAggregateWithTenantAndID(request.Tenant, opportunityId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateRenewalOpportunity) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created opportunity
	return &opportunitypb.OpportunityIdGrpcResponse{Id: opportunityId}, nil
}

func (s *opportunityService) UpdateRenewalOpportunity(ctx context.Context, request *opportunitypb.UpdateRenewalOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OpportunityService.UpdateRenewalOpportunity")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.Id)

	// Check if the opportunity ID is valid
	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	initAggregateFunc := func() eventstore.Aggregate {
		return aggregate.NewOpportunityAggregateWithTenantAndID(request.Tenant, request.Id)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateRenewalOpportunity.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &opportunitypb.OpportunityIdGrpcResponse{Id: request.Id}, nil
}

func (s *opportunityService) CloseLooseOpportunity(ctx context.Context, request *opportunitypb.CloseLooseOpportunityGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OpportunityService.CloseLooseOpportunity")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.Id)

	// Check if the opportunity ID is valid
	if request.Id == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("id"))
	}

	// Convert any protobuf timestamp to time.Time, if necessary
	closedAt := utils.TimestampProtoToTimePtr(request.ClosedAt)

	closeLooseOpportunityCommand := command.NewCloseLooseOpportunityCommand(
		request.Id,
		request.Tenant,
		request.LoggedInUserId,
		request.AppSource,
		nil,
		closedAt)

	if err := s.opportunityCommandHandlers.CloseLooseOpportunity.Handle(ctx, closeLooseOpportunityCommand); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CloseLooseOpportunity.Handle) tenant:{%v}, opportunityId:{%v}, err: %v", request.Tenant, request.Id, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created opportunity
	return &opportunitypb.OpportunityIdGrpcResponse{Id: request.Id}, nil
}

func (s *opportunityService) UpdateRenewalOpportunityNextCycleDate(ctx context.Context, request *opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest) (*opportunitypb.OpportunityIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OpportunityService.UpdateRenewalOpportunityNextCycleDate")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)
	span.SetTag(tracing.SpanTagEntityId, request.OpportunityId)

	// Check if the opportunity ID is valid
	if request.OpportunityId == "" {
		return nil, grpcerr.ErrResponse(grpcerr.ErrMissingField("opportunityId"))
	}

	cmd := command.NewUpdateRenewalOpportunityNextCycleDateCommand(
		request.OpportunityId,
		request.Tenant,
		request.LoggedInUserId,
		request.AppSource,
		nil,
		utils.TimestampProtoToTimePtr(request.RenewedAt))

	if err := s.opportunityCommandHandlers.UpdateRenewalOpportunityNextCycleDate.Handle(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateRenewalOpportunityNextCycleDate.Handle) tenant:{%v}, opportunityId:{%v}, err: %s", request.Tenant, request.OpportunityId, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	// Return the ID of the newly created opportunity
	return &opportunitypb.OpportunityIdGrpcResponse{Id: request.OpportunityId}, nil
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

func (s *opportunityService) checkContractExists(ctx context.Context, tenant, contractId string) (bool, error) {
	contractAggregate := contractaggregate.NewContractAggregateWithTenantAndID(tenant, contractId)
	err := s.aggregateStore.Exists(ctx, contractAggregate.GetID())
	if err != nil {
		if errors.Is(err, eventstore.ErrAggregateNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil // The contract exists
}

func containsOpportunityMaskFieldAll(fields []opportunitypb.OpportunityMaskField) bool {
	for _, field := range fields {
		if field == opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_ALL {
			return true
		}
	}
	return false
}
