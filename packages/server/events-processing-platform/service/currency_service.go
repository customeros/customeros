package service

import (
	"context"
	"github.com/google/uuid"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/currency"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	currencypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/currency"
)

type currencyService struct {
	currencypb.UnimplementedCurrencyGrpcServiceServer
	log           logger.Logger
	eventHandlers *currency.EventHandlers
}

func NewCurrencyService(log logger.Logger, eventHandlers *currency.EventHandlers) *currencyService {
	return &currencyService{
		log:           log,
		eventHandlers: eventHandlers,
	}
}

func (s *currencyService) CreateCurrency(ctx context.Context, request *currencypb.CreateCurrencyRequest) (*currencypb.CurrencyIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "CurrencyService.CreateCurrency")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, "", request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	currencyId := uuid.New().String()

	baseRequest := eventstore.NewBaseRequest(currencyId, "", request.LoggedInUserId, commonmodel.SourceFromGrpc(request.SourceFields))

	if err := s.eventHandlers.CurrencyCreate.Handle(ctx, baseRequest, request); err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("(CurrencyService.CurrencyCreate), err: %v", err.Error())
		return nil, grpcerr.ErrResponse(err)
	}

	return &currencypb.CurrencyIdGrpcResponse{Id: currencyId}, nil
}

func (s *currencyService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
