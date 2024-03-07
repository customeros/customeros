package service

import (
	"context"
	"github.com/google/uuid"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/country"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	countrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/country"
)

type countryService struct {
	countrypb.UnimplementedCountryGrpcServiceServer
	log           logger.Logger
	eventHandlers *country.EventHandlers
}

func NewCountryService(log logger.Logger, eventHandlers *country.EventHandlers) *countryService {
	return &countryService{
		log:           log,
		eventHandlers: eventHandlers,
	}
}

func (s *countryService) CreateCountry(ctx context.Context, request *countrypb.CreateCountryRequest) (*countrypb.CountryIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "CountryService.CreateCountry")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, "", request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	countryId := uuid.New().String()

	baseRequest := eventstore.NewBaseRequest(countryId, "", request.LoggedInUserId, commonmodel.SourceFromGrpc(request.SourceFields))

	if err := s.eventHandlers.CountryCreate.Handle(ctx, baseRequest, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CountryService.CountryCreate), err: %v", err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &countrypb.CountryIdGrpcResponse{Id: countryId}, nil
}

func (s *countryService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
