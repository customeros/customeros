package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/errors"
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

	organizationAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(command.Tenant, command.ObjectID)
	err := c.es.Exists(ctx, organizationAggregate.GetID())
	if err != nil {
		return eventstore.ErrInvalidAggregate
	} else {
		organizationAggregate, _ = aggregate.LoadOrganizationAggregate(ctx, c.es, command.Tenant, command.ObjectID)
		if err = organizationAggregate.LinkPhoneNumber(ctx, command.Tenant, command.PhoneNumberId, command.Label, command.Primary); err != nil {
			return err
		}
		if command.Primary {
			for k, v := range organizationAggregate.Organization.PhoneNumbers {
				if k != command.PhoneNumberId && v.Primary {
					if err = organizationAggregate.SetPhoneNumberNonPrimary(ctx, command.Tenant, command.PhoneNumberId); err != nil {
						return err
					}
				}
			}
		}
	}

	span.LogFields(log.String("Organization", organizationAggregate.Organization.String()))
	return c.es.Save(ctx, organizationAggregate)
}
