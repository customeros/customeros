package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type RequestSummaryCommandHandler interface {
	Handle(ctx context.Context, cmd *command.RequestSummaryCommand) error
}

type requestSummaryCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewRequestSummaryCommandHandler(log logger.Logger, es eventstore.AggregateStore) RequestSummaryCommandHandler {
	return &requestSummaryCommandHandler{log: log, es: es}
}

func (c *requestSummaryCommandHandler) Handle(ctx context.Context, cmd *command.RequestSummaryCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestSummaryCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", cmd.Tenant), log.String("ObjectID", cmd.ObjectID))

	if err := validator.GetValidator().Struct(cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	interactionEventAggregate, err := aggregate.LoadInteractionEventAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = interactionEventAggregate.RequestSummary(ctx, cmd.Tenant); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, interactionEventAggregate)
}
