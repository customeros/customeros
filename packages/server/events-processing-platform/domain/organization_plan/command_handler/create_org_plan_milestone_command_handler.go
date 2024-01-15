package command_handler

import (
	"context"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CreateOrganizationPlanMilestoneCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CreateOrganizationPlanMilestoneCommand) error
}

type createOrganizationPlanMilestoneCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewCreateOrganizationPlanMilestoneCommandHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) CreateOrganizationPlanMilestoneCommandHandler {
	return &createOrganizationPlanMilestoneCommandHandler{log: log, es: es, cfg: cfg}
}

// Handle processes the CreateOrganizationPlanCommand to create a new org plan.
func (h *createOrganizationPlanMilestoneCommandHandler) Handle(ctx context.Context, cmd *command.CreateOrganizationPlanMilestoneCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createOrganizationPlanMilestoneCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	tracing.LogObjectAsJson(span, "command", cmd)

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		// Load or initialize the org aggregate
		orgAggregate, err := aggregate.LoadOrganizationAggregate(ctx, h.es, cmd.Tenant, cmd.GetObjectID())
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if eventstore.IsAggregateNotFound(orgAggregate) {
			tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
			return eventstore.ErrAggregateNotFound
		}

		// Apply the command to the aggregate
		if err = orgAggregate.HandleCommand(ctx, cmd); err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		// Persist the changes to the event store
		err = h.es.Save(ctx, orgAggregate)
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
