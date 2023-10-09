package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type AddSocialCommandHandler interface {
	Handle(ctx context.Context, command *command.AddSocialCommand) error
}

type addSocialCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewAddSocialCommandHandler(log logger.Logger, es eventstore.AggregateStore) AddSocialCommandHandler {
	return &addSocialCommandHandler{log: log, es: es}
}

func (c *addSocialCommandHandler) Handle(ctx context.Context, command *command.AddSocialCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AddSocialCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, command.Tenant, command.UserID)
	span.LogFields(log.String("ObjectID", command.ObjectID))

	if command.Tenant == "" {
		return eventstore.ErrMissingTenant
	}

	organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if err = organizationAggregate.HandleCommand(ctx, command); err != nil {
		return err
	}

	return c.es.Save(ctx, organizationAggregate)
}
