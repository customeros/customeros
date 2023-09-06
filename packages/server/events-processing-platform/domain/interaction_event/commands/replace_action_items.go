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
	"time"
)

type ReplaceActionItemsCommand struct {
	eventstore.BaseCommand
	ActionItems []string
	UpdatedAt   *time.Time
}

func NewReplaceActionItemsCommand(tenant, interactionEventId string, actionItems []string, updatedAt *time.Time) *ReplaceActionItemsCommand {
	return &ReplaceActionItemsCommand{
		BaseCommand: eventstore.NewBaseCommand(interactionEventId, tenant, ""),
		ActionItems: actionItems,
		UpdatedAt:   updatedAt,
	}
}

type ReplaceActionItemsCommandHandler interface {
	Handle(ctx context.Context, command *ReplaceActionItemsCommand) error
}

type replaceActionItemsCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewReplaceActionItemsCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) ReplaceActionItemsCommandHandler {
	return &replaceActionItemsCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *replaceActionItemsCommandHandler) Handle(ctx context.Context, command *ReplaceActionItemsCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReplaceActionItemsCommandHandler.Handle")
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

	if err = interactionEventAggregate.ReplaceActionItems(ctx, command.Tenant, command.ActionItems, command.UpdatedAt); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, interactionEventAggregate)
}
