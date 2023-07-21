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
)

type UpsertUserCommandHandler interface {
	Handle(ctx context.Context, command *UpsertUserCommand) error
}

type upsertUserCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertUserCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) UpsertUserCommandHandler {
	return &upsertUserCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *upsertUserCommandHandler) Handle(ctx context.Context, command *UpsertUserCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertUserCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("UserID", command.ObjectID))

	if command.Tenant == "" {
		tracing.TraceErr(span, eventstore.ErrMissingTenant)
		return eventstore.ErrMissingTenant
	}

	userAggregate, err := aggregate.LoadUserAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if aggregate.IsAggregateNotFound(userAggregate) {
		if err = userAggregate.CreateUser(ctx, UpsertUserCommandToUserDto(command)); err != nil {
			return err
		}
	} else {
		if err = userAggregate.UpdateUser(ctx, UpsertUserCommandToUserDto(command)); err != nil {
			return err
		}
	}

	return c.es.Save(ctx, userAggregate)
}
