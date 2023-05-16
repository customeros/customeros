package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
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

func NewLinkPhoneNumberCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *linkPhoneNumberCommandHandler {
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

	contactAggregate := aggregate.NewContactAggregateWithTenantAndID(command.Tenant, command.ObjectID)
	err := c.es.Exists(ctx, contactAggregate.GetID())
	if err != nil {
		return eventstore.ErrInvalidAggregate
	} else {
		contactAggregate, _ = aggregate.LoadContactAggregate(ctx, c.es, command.Tenant, command.ObjectID)
		if err = contactAggregate.LinkPhoneNumber(ctx, command.Tenant, command.PhoneNumberId, command.Label, command.Primary); err != nil {
			return err
		}
		if command.Primary {
			for k, v := range contactAggregate.Contact.PhoneNumbers {
				if k != command.PhoneNumberId && v.Primary {
					if err = contactAggregate.SetPhoneNumberNonPrimary(ctx, command.Tenant, command.PhoneNumberId); err != nil {
						return err
					}
				}
			}
		}
	}

	span.LogFields(log.String("Contact", contactAggregate.Contact.String()))
	return c.es.Save(ctx, contactAggregate)
}
