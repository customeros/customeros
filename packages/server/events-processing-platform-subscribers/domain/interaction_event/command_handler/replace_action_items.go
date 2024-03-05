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

type ReplaceActionItemsCommandHandler interface {
	Handle(ctx context.Context, cmd *command.ReplaceActionItemsCommand) error
}

type replaceActionItemsCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewReplaceActionItemsCommandHandler(log logger.Logger, es eventstore.AggregateStore) ReplaceActionItemsCommandHandler {
	return &replaceActionItemsCommandHandler{log: log, es: es}
}

func (h *replaceActionItemsCommandHandler) Handle(ctx context.Context, cmd *command.ReplaceActionItemsCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReplaceActionItemsCommandHandler.Handle")
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

	if err = interactionEventAggregate.ReplaceActionItems(ctx, cmd.Tenant, cmd.ActionItems, cmd.UpdatedAt); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, interactionEventAggregate)
}
