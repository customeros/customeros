package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/country"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	countrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/country"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type countryService struct {
	countrypb.UnimplementedCountryGrpcServiceServer
	services       *Services
	log            logger.Logger
	aggregateStore eventstore.AggregateStore
	cfg            *config.Config
}

func NewCountryService(services *Services, log logger.Logger, aggregateStore eventstore.AggregateStore, cfg *config.Config) *countryService {
	return &countryService{
		services:       services,
		log:            log,
		aggregateStore: aggregateStore,
		cfg:            cfg,
	}
}

func (s *countryService) CreateCountry(ctx context.Context, request *countrypb.CreateCountryRequest) (*countrypb.CountryIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "CountryService.CreateCountry")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, "", request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	countryId := uuid.New().String()

	initAggregateFunc := func() eventstore.Aggregate {
		return country.NewCountryAggregateWithID(countryId)
	}
	if _, err := s.services.RequestHandler.HandleGRPCRequest(ctx, initAggregateFunc, eventstore.LoadAggregateOptions{}, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CountryService.CountryCreate), err: %v", err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &countrypb.CountryIdGrpcResponse{Id: countryId}, nil
}

func (s *countryService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
