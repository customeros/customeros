package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LinkPhoneNumberCommandHandler interface {
	Handle(ctx context.Context, command *LinkPhoneNumberCommand) error
}

type linkPhoneNumberCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewLinkPhoneNumberCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) LinkPhoneNumberCommandHandler {
	return &linkPhoneNumberCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *linkPhoneNumberCommandHandler) Handle(ctx context.Context, command *LinkPhoneNumberCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "linkPhoneNumberCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if len(command.Tenant) == 0 {
		return eventstore.ErrMissingTenant
	}
	if len(command.PhoneNumberId) == 0 {
		return errors.ErrMissingPhoneNumberId
	}

	userAggregate, err := aggregate.LoadUserAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = userAggregate.LinkPhoneNumber(ctx, command.Tenant, command.PhoneNumberId, command.Label, command.Primary); err != nil {
		return err
	}
	if command.Primary {
		for k, v := range userAggregate.User.PhoneNumbers {
			if k != command.PhoneNumberId && v.Primary {
				if err = userAggregate.SetPhoneNumberNonPrimary(ctx, command.Tenant, command.PhoneNumberId); err != nil {
					return err
				}
			}
		}
	}

	span.LogFields(log.String("User", userAggregate.User.String()))
	return c.es.Save(ctx, userAggregate)
}
