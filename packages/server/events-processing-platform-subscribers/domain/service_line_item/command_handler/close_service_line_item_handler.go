package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

// CloseServiceLineItemCommandHandler defines the interface for a handler that can process CloseServiceLineItemCommands.
type CloseServiceLineItemCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CloseServiceLineItemCommand) error
}

type closeServiceLineItemCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

// NewCloseServiceLineItemCommandHandler creates a new handler for updating service line items.
func NewCloseServiceLineItemCommandHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) CloseServiceLineItemCommandHandler {
	return &closeServiceLineItemCommandHandler{log: log, es: es, cfg: cfg}
}

// Handle processes the CloseServiceLineItemCommand to close an existing service line item.
func (h *closeServiceLineItemCommandHandler) Handle(ctx context.Context, cmd *command.CloseServiceLineItemCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "closeServiceLineItemCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		// Load the service line item aggregate
		sliAggregate, err := aggregate.LoadServiceLineItemAggregate(ctx, h.es, cmd.Tenant, cmd.GetObjectID())
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		// Apply the command to the aggregate
		if err = sliAggregate.HandleCommand(ctx, cmd); err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		// Persist the changes to the event store
		err = h.es.Save(ctx, sliAggregate)
		if err == nil {
			return nil // Save successful
		}

		if eventstore.IsEventStoreErrorCodeWrongExpectedVersion(err) {
			// Handle concurrency error
			if attempt == h.cfg.RetriesOnOptimisticLockException-1 {
				// If we have reached the maximum number of retries, return an error
				tracing.TraceErr(span, err)
				return err
			}
			span.LogFields(log.Int("retryAttempt", attempt+1))
			time.Sleep(utils.BackOffExponentialDelay(attempt)) // backoffDelay is a function that increases the delay with each attempt
			continue                                           // Retry
		} else {
			// Some other error occurred
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
