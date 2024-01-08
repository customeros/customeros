package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
)

type CreateInvoicingCycleCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CreateInvoicingCycleTypeCommand) error
}

type createInvoicingCycleCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewCreateInvoicingCycleCommandHandler(log logger.Logger, es eventstore.AggregateStore) CreateInvoicingCycleCommandHandler {
	return &createInvoicingCycleCommandHandler{log: log, es: es}
}

func (h *createInvoicingCycleCommandHandler) Handle(ctx context.Context, cmd *command.CreateInvoicingCycleTypeCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateInvoicingCycleCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	tracing.LogObjectAsJson(span, "command", cmd)

	invoicingCycleAggregate, err := aggregate.LoadInvoicingCycleAggregate(ctx, h.es, cmd.Tenant, cmd.GetObjectID())
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Apply the command to the aggregate
	if err = invoicingCycleAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err = h.es.Save(ctx, invoicingCycleAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
