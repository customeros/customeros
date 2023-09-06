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

type RequestSummaryCommand struct {
	eventstore.BaseCommand
}

func NewRequestSummaryCommand(tenant, interactionEventId string) *RequestSummaryCommand {
	return &RequestSummaryCommand{
		BaseCommand: eventstore.NewBaseCommand(interactionEventId, tenant, ""),
	}
}

type RequestSummaryCommandHandler interface {
	Handle(ctx context.Context, command *RequestSummaryCommand) error
}

type requestSummaryCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewRequestSummaryCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) RequestSummaryCommandHandler {
	return &requestSummaryCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *requestSummaryCommandHandler) Handle(ctx context.Context, command *RequestSummaryCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestSummaryCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if err := validator.GetValidator().Struct(command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	interactionEventAggregate, err := aggregate.LoadInteractionEventAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = interactionEventAggregate.RequestSummary(ctx, command.Tenant); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, interactionEventAggregate)
}
