package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LocationValidatedCommandHandler interface {
	Handle(ctx context.Context, cmd *command.LocationValidatedCommand) error
}

type locationValidatedCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewLocationValidatedCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) LocationValidatedCommandHandler {
	return &locationValidatedCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *locationValidatedCommandHandler) Handle(ctx context.Context, cmd *command.LocationValidatedCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationValidatedCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", cmd.Tenant), log.String("ObjectID", cmd.ObjectID))

	locationAggregate, err := aggregate.LoadLocationAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = locationAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, locationAggregate)
}
