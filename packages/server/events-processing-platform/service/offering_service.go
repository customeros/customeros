package service

import (
	"context"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"

	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/offering"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	offeringpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/offering"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type offeringService struct {
	offeringpb.UnimplementedOfferingGrpcServiceServer
	log            logger.Logger
	requestHandler offering.OfferingRequestHandler
}

func NewOfferingService(log logger.Logger, aggregateStore eventstore.AggregateStore, cfg *config.Config) *offeringService {
	return &offeringService{
		log:            log,
		requestHandler: offering.NewOfferingRequestHandler(log, aggregateStore, cfg.Utils),
	}
}

func (s *offeringService) CreateOffering(ctx context.Context, request *offeringpb.CreateOfferingGrpcRequest) (*commonpb.IdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OfferingService.CreateOffering")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	offeringId := uuid.New().String()

	_, err := s.requestHandler.Handle(ctx, request.Tenant, offeringId, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CreateOffering.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &commonpb.IdResponse{Id: offeringId}, nil
}

func (s *offeringService) UpdateOffering(ctx context.Context, request *offeringpb.UpdateOfferingGrpcRequest) (*commonpb.IdResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OfferingService.UpdateOffering")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.Id)
	tracing.LogObjectAsJson(span, "request", request)

	_, err := s.requestHandler.HandleWithRetry(ctx, request.Tenant, request.Id, false, request)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(UpdateOffering.Handle) tenant:{%v}, err: %v", request.Tenant, err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &commonpb.IdResponse{Id: request.Id}, nil
}
