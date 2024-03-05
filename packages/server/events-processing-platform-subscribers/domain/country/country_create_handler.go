package country

import (
	"context"
	commonAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	countrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/country"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"time"
)

type CountryCreateHandler interface {
	Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *countrypb.CreateCountryRequest) error
}

type countryCreateHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewCountryCreateHandler(log logger.Logger, es eventstore.AggregateStore) CountryCreateHandler {
	return &countryCreateHandler{log: log, es: es}
}

func (h *countryCreateHandler) Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *countrypb.CreateCountryRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryCreateHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, baseRequest.Tenant, baseRequest.LoggedInUserId)
	tracing.LogObjectAsJson(span, "common", baseRequest)
	tracing.LogObjectAsJson(span, "request", request)

	countryAggregate, err := LoadCountryAggregate(ctx, h.es, baseRequest.Tenant, baseRequest.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	createEvent, err := NewCountryCreateEvent(countryAggregate, request.Name, request.CodeA2, request.CodeA3, request.PhoneCode, time.Now().UTC(), baseRequest.SourceFields)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "CountryCreateEvent")
	}
	commonAggregate.EnrichEventWithMetadataExtended(&createEvent, span, commonAggregate.EventMetadata{
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	err = countryAggregate.Apply(createEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err = h.es.Save(ctx, countryAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
