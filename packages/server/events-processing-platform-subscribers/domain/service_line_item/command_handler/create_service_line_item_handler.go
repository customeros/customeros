package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// CreateServiceLineItemCommandHandler defines the interface for a handler that can process CreateServiceLineItemCommands.
type CreateServiceLineItemCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CreateServiceLineItemCommand) error
}

type createServiceLineItemCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

// NewCreateServiceLineItemCommandHandler creates a new handler for creating service line items.
func NewCreateServiceLineItemCommandHandler(log logger.Logger, es eventstore.AggregateStore) CreateServiceLineItemCommandHandler {
	return &createServiceLineItemCommandHandler{log: log, es: es}
}

// Handle processes the CreateServiceLineItemCommand to create a new service line item.
func (h *createServiceLineItemCommandHandler) Handle(ctx context.Context, cmd *command.CreateServiceLineItemCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createServiceLineItemCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	// Load or initialize the service line item aggregate
	serviceLineItemAggregate, err := aggregate.LoadServiceLineItemAggregate(ctx, h.es, cmd.Tenant, cmd.GetObjectID())
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Apply the command to the aggregate
	if err = serviceLineItemAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err = h.es.Save(ctx, serviceLineItemAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
