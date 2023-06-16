package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LinkJobRoleCommandHandler interface {
	Handle(ctx context.Context, command *LinkJobRoleCommand) error
}

type linkJobRoleCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewLinkJobRoleCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *linkJobRoleCommandHandler {
	return &linkJobRoleCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *linkJobRoleCommandHandler) Handle(ctx context.Context, command *LinkJobRoleCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "linkJobRoleCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if len(command.Tenant) == 0 {
		return eventstore.ErrMissingTenant
	}
	if len(command.JobRoleId) == 0 {
		return errors.ErrMissingEmailId
	}

	userAggregate := aggregate.NewUserAggregateWithTenantAndID(command.Tenant, command.ObjectID)
	err := c.es.Exists(ctx, userAggregate.GetID())
	if err != nil {
		return eventstore.ErrInvalidAggregate
	} else {
		userAggregate, _ = aggregate.LoadUserAggregate(ctx, c.es, command.Tenant, command.ObjectID)
		if err = userAggregate.LinkJobRole(ctx, command.Tenant, command.JobRoleId); err != nil {
			return err
		}
	}

	span.LogFields(log.String("User", userAggregate.User.String()))
	return c.es.Save(ctx, userAggregate)
}
