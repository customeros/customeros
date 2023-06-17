package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
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
	span.LogFields(log.String("Tenant", command.Tenant), log.String("UserID", command.ObjectID))

	if len(command.Tenant) == 0 {
		tracing.TraceErr(span, eventstore.ErrMissingTenant)
		return eventstore.ErrMissingTenant
	}

	userAggregate := aggregate.NewUserAggregateWithTenantAndID(command.Tenant, command.ObjectID)
	err := c.es.Exists(ctx, userAggregate.GetID())
	if err != nil && !errors.Is(err, eventstore.ErrAggregateNotFound) {
		return err
	} else if err != nil && errors.Is(err, eventstore.ErrAggregateNotFound) {
		if err = userAggregate.CreateUser(ctx, UpsertUserCommandToUserDto(command)); err != nil {
			return err
		}
	} else {
		userAggregate, _ = aggregate.LoadUserAggregate(ctx, c.es, command.Tenant, command.ObjectID)
		if err = userAggregate.UpdateUser(ctx, UpsertUserCommandToUserDto(command)); err != nil {
			return err
		}
	}

	span.LogFields(log.String("User", userAggregate.User.String()))
	return c.es.Save(ctx, userAggregate)
}
