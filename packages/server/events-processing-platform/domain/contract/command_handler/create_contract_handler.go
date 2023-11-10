package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// CreateContractCommandHandler defines the interface for a handler that can process CreateContractCommands.
type CreateContractCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CreateContractCommand) error
}

type createContractCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

// NewCreateContractCommandHandler creates a new handler for creating contracts.
func NewCreateContractCommandHandler(log logger.Logger, es eventstore.AggregateStore) CreateContractCommandHandler {
	return &createContractCommandHandler{log: log, es: es}
}

// Handle processes the CreateContractCommand to create a new contract.
func (h *createContractCommandHandler) Handle(ctx context.Context, cmd *command.CreateContractCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createContractCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	// Validate the command fields
	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	// Load or initialize the contract aggregate
	contractAggregate, err := aggregate.LoadContractAggregate(ctx, h.es, cmd.Tenant, cmd.GetObjectID())
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Apply the command to the aggregate
	if err := contractAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err := h.es.Save(ctx, contractAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
