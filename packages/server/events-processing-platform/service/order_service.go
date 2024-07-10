package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/order"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	orderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/order"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderService struct {
	orderpb.UnimplementedOrderGrpcServiceServer
	log                 logger.Logger
	orderRequestHandler order.OrderRequestHandler
	aggregateStore      eventstore.AggregateStore
}

func NewOrderService(log logger.Logger, aggregateStore eventstore.AggregateStore, cfg *config.Config) *orderService {
	return &orderService{
		log:                 log,
		orderRequestHandler: order.NewOrderRequestHandler(log, aggregateStore, cfg.Utils),
		aggregateStore:      aggregateStore,
	}
}

func (s *orderService) UpsertOrder(ctx context.Context, request *orderpb.UpsertOrderGrpcRequest) (*orderpb.OrderIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrderService.UpsertOrder")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	// Check if the contract aggregate exists
	organizationExists, err := s.checkOrganizationExists(ctx, request.Tenant, request.OrganizationId)
	if err != nil {
		s.log.Error(err, "error checking contract existence")
		return nil, status.Errorf(codes.Internal, "error checking contract existence: %v", err)
	}
	if !organizationExists {
		return nil, status.Errorf(codes.NotFound, "organization with ID %s not found", request.OrganizationId)
	}

	orderId := utils.FirstNotEmpty(request.Id, uuid.New().String())

	if _, err := s.orderRequestHandler.HandleWithRetry(ctx, request.Tenant, *orderId, false, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpsertOrder) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &orderpb.OrderIdGrpcResponse{Id: *orderId}, nil
}

func (s *orderService) checkOrganizationExists(ctx context.Context, tenant, organizationId string) (bool, error) {
	organizationAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(tenant, organizationId)
	err := s.aggregateStore.Exists(ctx, organizationAggregate.GetID())
	if err != nil {
		if errors.Is(err, eventstore.ErrAggregateNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
