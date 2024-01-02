package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
)

type CreateMasterPlanCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CreateMasterPlanCommand) error
}

type createMasterPlanCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewCreateMasterPlanCommandHandler(log logger.Logger, es eventstore.AggregateStore) CreateMasterPlanCommandHandler {
	return &createMasterPlanCommandHandler{log: log, es: es}
}

// Handle processes the CreateMasterPlanCommand to create a new master plan.
func (h *createMasterPlanCommandHandler) Handle(ctx context.Context, cmd *command.CreateMasterPlanCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createMasterPlanCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	tracing.LogObjectAsJson(span, "command", cmd)

	// Load or initialize the master plan aggregate
	masterPlanAggregate, err := aggregate.LoadMasterPlanAggregate(ctx, h.es, cmd.Tenant, cmd.GetObjectID())
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Apply the command to the aggregate
	if err = masterPlanAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err = h.es.Save(ctx, masterPlanAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
