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

func (c *upsertCustomFieldCommandHandler) Handle(ctx context.Context, command *command.UpsertCustomFieldCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpsertCustomFieldCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, command.Tenant, command.LoggedInUserId)
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if err := validator.GetValidator().Struct(command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if aggregate.IsAggregateNotFound(organizationAggregate) {
		tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
		return eventstore.ErrAggregateNotFound
	} else {
		if err = organizationAggregate.HandleCommand(ctx, command); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return c.es.Save(ctx, organizationAggregate)
}
