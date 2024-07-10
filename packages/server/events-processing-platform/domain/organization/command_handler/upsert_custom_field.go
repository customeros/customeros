package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
)

type UpsertCustomFieldCommandHandler interface {
	Handle(ctx context.Context, command *command.UpsertCustomFieldCommand) error
}

type upsertCustomFieldCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewUpsertCustomFieldCommandHandler(log logger.Logger, es eventstore.AggregateStore) UpsertCustomFieldCommandHandler {
	return &upsertCustomFieldCommandHandler{log: log, es: es}
}

func (c *upsertCustomFieldCommandHandler) Handle(ctx context.Context, cmd *command.UpsertCustomFieldCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpsertCustomFieldCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	tracing.LogObjectAsJson(span, "command", cmd)

	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		return err
	}

	if eventstore.IsAggregateNotFound(organizationAggregate) {
		tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
		return eventstore.ErrAggregateNotFound
	} else {
		if err = organizationAggregate.HandleCommand(ctx, cmd); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return c.es.Save(ctx, organizationAggregate)
}
