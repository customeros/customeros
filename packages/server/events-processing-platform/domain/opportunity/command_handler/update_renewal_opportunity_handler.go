package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// UpdateRenewalOpportunityCommandHandler defines the interface for a handler that can process UpdateRenewalOpportunityCommands.
type UpdateRenewalOpportunityCommandHandler interface {
	Handle(ctx context.Context, cmd *command.UpdateRenewalOpportunityCommand) error
}

type updateRenewalOpportunityCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

// NewUpdateRenewalOpportunityCommandHandler updates a new handler for creating renewal opportunities.
func NewUpdateRenewalOpportunityCommandHandler(log logger.Logger, es eventstore.AggregateStore) UpdateRenewalOpportunityCommandHandler {
	return &updateRenewalOpportunityCommandHandler{
		log: log,
		es:  es,
	}
}

// Handle processes the UpdateRenewalOpportunityCommand to update a new renewal opportunity.
func (h *updateRenewalOpportunityCommandHandler) Handle(ctx context.Context, cmd *command.UpdateRenewalOpportunityCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateRenewalOpportunityCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	// Validate the command fields
	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	// Initialize the opportunity aggregate
	opportunityAggregate, err := aggregate.LoadOpportunityAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if aggregate.IsAggregateNotFound(opportunityAggregate) {
		tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
		return eventstore.ErrAggregateNotFound
	}

	// Apply the command to the aggregate
	if err = opportunityAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err = h.es.Save(ctx, opportunityAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
