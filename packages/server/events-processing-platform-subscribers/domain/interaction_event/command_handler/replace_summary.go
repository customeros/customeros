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

type ReplaceSummaryCommandHandler interface {
	Handle(ctx context.Context, cmd *command.ReplaceSummaryCommand) error
}

type replaceSummaryCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewReplaceSummaryCommandHandler(log logger.Logger, es eventstore.AggregateStore) ReplaceSummaryCommandHandler {
	return &replaceSummaryCommandHandler{log: log, es: es}
}

func (h *replaceSummaryCommandHandler) Handle(ctx context.Context, cmd *command.ReplaceSummaryCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReplaceSummaryCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", cmd.Tenant), log.String("ObjectID", cmd.ObjectID))

	if err := validator.GetValidator().Struct(cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	interactionEventAggregate, err := aggregate.LoadInteractionEventAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = interactionEventAggregate.ReplaceSummary(ctx, cmd.Tenant, cmd.Summary, cmd.ContentType, cmd.UpdatedAt); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return h.es.Save(ctx, interactionEventAggregate)
}
