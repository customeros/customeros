package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpsertLocationCommandHandler interface {
	Handle(ctx context.Context, cmd *command.UpsertLocationCommand) error
}

type upsertLocationHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertLocationHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) UpsertLocationCommandHandler {
	return &upsertLocationHandler{log: log, cfg: cfg, es: es}
}

func (h *upsertLocationHandler) Handle(ctx context.Context, cmd *command.UpsertLocationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpsertLocationHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", cmd.Tenant), log.String("ObjectID", cmd.ObjectID))

	if err := validator.GetValidator().Struct(cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	locationAggregate, err := aggregate.LoadLocationAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if aggregate.IsAggregateNotFound(locationAggregate) {
		cmd.IsCreateCommand = true
	}

	if err = locationAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	span.LogFields(log.String("Location aggregate", locationAggregate.Location.String()))
	return h.es.Save(ctx, locationAggregate)
}
