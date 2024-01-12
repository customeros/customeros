package command_handler

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
)

type CreateOrgPlanCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CreateOrgPlanCommand) error
}

type createOrgPlanCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewCreateOrgPlanCommandHandler(log logger.Logger, es eventstore.AggregateStore) CreateOrgPlanCommandHandler {
	return &createOrgPlanCommandHandler{log: log, es: es}
}

// Handle processes the CreateOrgPlanCommand to create a new master plan.
func (h *createOrgPlanCommandHandler) Handle(ctx context.Context, cmd *command.CreateOrgPlanCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createOrgPlanCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	tracing.LogObjectAsJson(span, "command", cmd)

	// Load or initialize the org plan aggregate
	orgPlanAggregate, err := aggregate.LoadOrganizationAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Apply the command to the aggregate
	if err = orgPlanAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err = h.es.Save(ctx, orgPlanAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
