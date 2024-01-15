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

type CreateOrganizationPlanCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CreateOrganizationPlanCommand) error
}

type createOrganizationPlanCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewCreateOrganizationPlanCommandHandler(log logger.Logger, es eventstore.AggregateStore) CreateOrganizationPlanCommandHandler {
	return &createOrganizationPlanCommandHandler{log: log, es: es}
}

// Handle processes the CreateOrganizationPlanCommand to create a new master plan.
func (h *createOrganizationPlanCommandHandler) Handle(ctx context.Context, cmd *command.CreateOrganizationPlanCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createOrganizationPlanCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	tracing.LogObjectAsJson(span, "command", cmd)

	// Load or initialize the org plan aggregate
	organizationPlanAggregate, err := aggregate.LoadOrganizationAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Apply the command to the aggregate
	if err = organizationPlanAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err = h.es.Save(ctx, organizationPlanAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
