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

type LinkEmailCommandHandler interface {
	Handle(ctx context.Context, command *LinkEmailCommand) error
}

type linkEmailCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewLinkEmailCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *linkEmailCommandHandler {
	return &linkEmailCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *linkEmailCommandHandler) Handle(ctx context.Context, command *LinkEmailCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "linkEmailCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if len(command.Tenant) == 0 {
		return eventstore.ErrMissingTenant
	}
	if len(command.EmailId) == 0 {
		return errors.ErrMissingEmailId
	}

	userAggregate := aggregate.NewUserAggregateWithTenantAndID(command.Tenant, command.ObjectID)
	err := c.es.Exists(ctx, userAggregate.GetID())
	if err != nil {
		return eventstore.ErrInvalidAggregate
	} else {
		userAggregate, _ = aggregate.LoadUserAggregate(ctx, c.es, command.Tenant, command.ObjectID)
		if err = userAggregate.LinkEmail(ctx, command.Tenant, command.EmailId, command.Label, command.Primary); err != nil {
			return err
		}
		if command.Primary {
			for k, v := range userAggregate.User.Emails {
				if k != command.EmailId && v.Primary {
					if err = userAggregate.SetEmailNonPrimary(ctx, command.Tenant, command.EmailId); err != nil {
						return err
					}
				}
			}
		}
	}

	span.LogFields(log.String("User", userAggregate.User.String()))
	return c.es.Save(ctx, userAggregate)
}
