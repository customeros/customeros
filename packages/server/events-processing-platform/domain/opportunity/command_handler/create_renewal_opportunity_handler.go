package command_handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

// CreateRenewalOpportunityCommandHandler defines the interface for a handler that can process CreateRenewalOpportunityCommands.
type CreateRenewalOpportunityCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CreateRenewalOpportunityCommand) error
}

type createRenewalOpportunityCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

// NewCreateRenewalOpportunityCommandHandler creates a new handler for creating renewal opportunities.
func NewCreateRenewalOpportunityCommandHandler(log logger.Logger, es eventstore.AggregateStore) CreateRenewalOpportunityCommandHandler {
	return &createRenewalOpportunityCommandHandler{
		log: log,
		es:  es,
	}
}

// Handle processes the CreateRenewalOpportunityCommand to create a new renewal opportunity.
func (h *createRenewalOpportunityCommandHandler) Handle(ctx context.Context, cmd *command.CreateRenewalOpportunityCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createRenewalOpportunityCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	if strings.TrimSpace(cmd.ObjectID) == "" {
		cmd.ObjectID = uuid.New().String()
	}

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
