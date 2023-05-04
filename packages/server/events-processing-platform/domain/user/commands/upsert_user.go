package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpsertUserCommandHandler interface {
	Handle(ctx context.Context, command *UpsertUserCommand) error
}

type upsertUserCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertUserCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *upsertUserCommandHandler {
	return &upsertUserCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *upsertUserCommandHandler) Handle(ctx context.Context, command *UpsertUserCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertUserCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", command.GetAggregateID()))

	if len(command.Tenant) == 0 {
		return eventstore.ErrMissingTenant
	}

	userAggregate := aggregate.NewUserAggregateWithTenantAndID(command.Tenant, command.AggregateID)
	err := c.es.Exists(ctx, userAggregate.GetID())
	if err != nil && !eventstore.IsErrEsResourceNotFound(err) {
		return err
	} else if err != nil && eventstore.IsErrEsResourceNotFound(err) {
		if err = userAggregate.CreateUser(ctx, UpsertUserCommandToUserDto(command)); err != nil {
			return err
		}
	} else {
		userAggregate, _ = aggregate.LoadUserAggregate(ctx, c.es, command.Tenant, command.AggregateID)
		if err = userAggregate.UpdateUser(ctx, UpsertUserCommandToUserDto(command)); err != nil {
			return err
		}
	}

	span.LogFields(log.String("User", userAggregate.User.String()))
	return c.es.Save(ctx, userAggregate)
}
