package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type RequestNextCycleDateCommandHandler interface {
	Handle(ctx context.Context, cmd *command.RequestNextCycleDateCommand) error
}

type requestNextCycleDateCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewRequestNextCycleDateCommandHandler(log logger.Logger, es eventstore.AggregateStore) RequestNextCycleDateCommandHandler {
	return &requestNextCycleDateCommandHandler{log: log, es: es}
}

// Handle processes the UpdateContractCommand to update a contract.
func (h *requestNextCycleDateCommandHandler) Handle(ctx context.Context, cmd *command.RequestNextCycleDateCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "requestNextCycleDateCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	// Load or initialize the contract aggregate
	contractAggregate, err := aggregate.LoadContractAggregate(ctx, h.es, cmd.Tenant, cmd.GetObjectID())
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Apply the command to the aggregate
	if err = contractAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err = h.es.Save(ctx, contractAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
