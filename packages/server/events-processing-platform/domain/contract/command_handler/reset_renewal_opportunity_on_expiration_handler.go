package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type ResetRenewalOpportunityOnExpirationCommandHandler interface {
	Handle(ctx context.Context, cmd *command.ResetRenewalOpportunityOnExpirationCommand) error
}

type resetRenewalOpportunityOnExpirationCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewResetRenewalOpportunityOnExpirationCommandHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) ResetRenewalOpportunityOnExpirationCommandHandler {
	return &resetRenewalOpportunityOnExpirationCommandHandler{
		log: log,
		es:  es,
		cfg: cfg,
	}
}

func (h *resetRenewalOpportunityOnExpirationCommandHandler) Handle(ctx context.Context, cmd *command.ResetRenewalOpportunityOnExpirationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ResetRenewalOpportunityOnExpirationCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
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
		err = h.es.Save(ctx, contractAggregate)
		if err == nil {
			return nil // Save successful
		}

		if eventstore.IsEventStoreErrorCodeWrongExpectedVersion(err) {
			// Handle concurrency error
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
