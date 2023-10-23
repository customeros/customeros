package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type RefreshLastTouchpointCommandHandler interface {
	Handle(ctx context.Context, command *command.RefreshLastTouchpointCommand) error
}

type refreshLastTouchpointCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewRefreshLastTouchpointCommandHandler(log logger.Logger, es eventstore.AggregateStore) RefreshLastTouchpointCommandHandler {
	return &refreshLastTouchpointCommandHandler{log: log, es: es}
}

func (c *refreshLastTouchpointCommandHandler) Handle(ctx context.Context, cmd *command.RefreshLastTouchpointCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RefreshLastTouchpointCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.String("ObjectID", cmd.ObjectID))

	if err := validator.GetValidator().Struct(cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if err = organizationAggregate.HandleCommand(ctx, cmd); err != nil {
		return err
	}

	return c.es.Save(ctx, organizationAggregate)
}
