package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type RequestActionItemsCommandHandler interface {
	Handle(ctx context.Context, cmd *command.RequestActionItemsCommand) error
}

type requestActionItemsCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewRequestActionItemsCommandHandler(log logger.Logger, es eventstore.AggregateStore) RequestActionItemsCommandHandler {
	return &requestActionItemsCommandHandler{log: log, es: es}
}

func (h *requestActionItemsCommandHandler) Handle(ctx context.Context, cmd *command.RequestActionItemsCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestActionItemsCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", cmd.Tenant), log.String("InteractionEventId", cmd.ObjectID))

	if err := validator.GetValidator().Struct(cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	interactionEventAggregate, err := aggregate.LoadInteractionEventAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = interactionEventAggregate.RequestActionItems(ctx, cmd.Tenant); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, interactionEventAggregate)
}
