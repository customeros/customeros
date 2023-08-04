package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type RequestActionItemsCommand struct {
	eventstore.BaseCommand
}

func NewRequestActionItemsCommand(tenant, interactionEventId string) *RequestActionItemsCommand {
	return &RequestActionItemsCommand{
		BaseCommand: eventstore.NewBaseCommand(interactionEventId, tenant),
	}
}

type RequestActionItemsCommandHandler interface {
	Handle(ctx context.Context, command *RequestActionItemsCommand) error
}

type requestActionItemsCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewRequestActionItemsCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) RequestActionItemsCommandHandler {
	return &requestActionItemsCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *requestActionItemsCommandHandler) Handle(ctx context.Context, command *RequestActionItemsCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestActionItemsCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("InteractionEventId", command.ObjectID))

	if err := validator.GetValidator().Struct(command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	interactionEventAggregate, err := aggregate.LoadInteractionEventAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = interactionEventAggregate.RequestActionItems(ctx, command.Tenant); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, interactionEventAggregate)
}
