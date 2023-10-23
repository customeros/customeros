package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LinkJobRoleCommandHandler interface {
	Handle(ctx context.Context, command *command.LinkJobRoleCommand) error
}

type linkJobRoleCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewLinkJobRoleCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) LinkJobRoleCommandHandler {
	return &linkJobRoleCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *linkJobRoleCommandHandler) Handle(ctx context.Context, command *command.LinkJobRoleCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "linkJobRoleCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if len(command.Tenant) == 0 {
		return eventstore.ErrMissingTenant
	}
	if len(command.JobRoleId) == 0 {
		return errors.ErrMissingEmailId
	}

	userAggregate, err := aggregate.LoadUserAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = userAggregate.LinkJobRole(ctx, command.Tenant, command.JobRoleId, command.LoggedInUserId); err != nil {
		return err
	}

	span.LogFields(log.String("User", userAggregate.User.String()))
	return c.es.Save(ctx, userAggregate)
}
