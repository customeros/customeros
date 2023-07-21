package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LocationValidatedCommandHandler interface {
	Handle(ctx context.Context, command *LocationValidatedCommand) error
}

type locationValidatedCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewLocationValidatedCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) LocationValidatedCommandHandler {
	return &locationValidatedCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *locationValidatedCommandHandler) Handle(ctx context.Context, command *LocationValidatedCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationValidatedCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	locationAggregate, err := aggregate.LoadLocationAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	locationAddress := models.LocationAddress{}
	locationAddress.From(command.LocationAddressFields)
	if err = locationAggregate.LocationValidated(ctx, command.Tenant, command.RawAddress, command.CountryForValidation, locationAddress); err != nil {
		return err
	}

	return h.es.Save(ctx, locationAggregate)
}
