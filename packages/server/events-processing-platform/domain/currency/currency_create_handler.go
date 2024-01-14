package currency

import (
	"context"
	commonAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	currencypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/currency"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"time"
)

type CurrencyCreateHandler interface {
	Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *currencypb.CreateCurrencyRequest) error
}

type currencyCreateHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewCurrencyCreateHandler(log logger.Logger, es eventstore.AggregateStore) CurrencyCreateHandler {
	return &currencyCreateHandler{log: log, es: es}
}

func (h *currencyCreateHandler) Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *currencypb.CreateCurrencyRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CurrencyCreateHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, baseRequest.Tenant, baseRequest.LoggedInUserId)
	tracing.LogObjectAsJson(span, "common", baseRequest)
	tracing.LogObjectAsJson(span, "request", request)

	currencyAggregate, err := LoadCurrencyAggregate(ctx, h.es, baseRequest.Tenant, baseRequest.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	createEvent, err := NewCurrencyCreateEvent(currencyAggregate, request.Name, request.Symbol, time.Now().UTC(), baseRequest.SourceFields)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "CurrencyCreateEvent")
	}
	commonAggregate.EnrichEventWithMetadataExtended(&createEvent, span, commonAggregate.EventMetadata{
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	err = currencyAggregate.Apply(createEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err = h.es.Save(ctx, currencyAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
