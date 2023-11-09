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

// CreateOpportunityCommandHandler defines the interface for a handler that can process CreateOpportunityCommands.
type CreateOpportunityCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CreateOpportunityCommand) error
}

type createOpportunityCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

// NewCreateOpportunityCommandHandler creates a new handler for creating opportunities.
func NewCreateOpportunityCommandHandler(log logger.Logger, es eventstore.AggregateStore) CreateOpportunityCommandHandler {
	return &createOpportunityCommandHandler{log: log, es: es}
}

// Handle processes the CreateOpportunityCommand to create a new opportunity.
func (h *createOpportunityCommandHandler) Handle(ctx context.Context, cmd *command.CreateOpportunityCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createOpportunityCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	// Validate the command fields
	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	// Load or initialize the opportunity aggregate
	opportunityAggregate, err := aggregate.LoadOpportunityAggregate(ctx, h.es, cmd.Tenant, cmd.GetObjectID())
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Apply the command to the aggregate
	if err := opportunityAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err := h.es.Save(ctx, opportunityAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
