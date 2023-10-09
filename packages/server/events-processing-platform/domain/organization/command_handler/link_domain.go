package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LinkDomainCommandHandler interface {
	Handle(ctx context.Context, command *command.LinkDomainCommand) error
}

type linkDomainCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewLinkDomainCommandHandler(log logger.Logger, es eventstore.AggregateStore) LinkDomainCommandHandler {
	return &linkDomainCommandHandler{log: log, es: es}
}

func (c *linkDomainCommandHandler) Handle(ctx context.Context, command *command.LinkDomainCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LinkDomainCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, command.Tenant, command.UserID)
	span.LogFields(log.String("ObjectID", command.ObjectID))

	if command.Tenant == "" {
		return eventstore.ErrMissingTenant
	}
	if command.Domain == "" {
		return errors.ErrMissingDomain
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
