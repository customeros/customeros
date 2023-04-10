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
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", command.GetAggregateID()))

	if len(command.Tenant) == 0 {
		return eventstore.ErrMissingTenant
	}
	if len(command.EmailId) == 0 {
		return errors.ErrMissingEmailId
	}

	organizationAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(command.Tenant, command.AggregateID)
	err := c.es.Exists(ctx, organizationAggregate.GetID())
	if err != nil {
		return eventstore.ErrInvalidAggregate
	} else {
		organizationAggregate, _ = aggregate.LoadOrganizationAggregate(ctx, c.es, command.Tenant, command.AggregateID)
		if err = organizationAggregate.LinkEmail(ctx, command.Tenant, command.EmailId, command.Label, command.Primary); err != nil {
			return err
		}
		if command.Primary {
			for k, v := range organizationAggregate.Organization.Emails {
				if k != command.EmailId && v.Primary {
					if err = organizationAggregate.SetEmailNonPrimary(ctx, command.Tenant, command.EmailId); err != nil {
						return err
					}
				}
			}
		}
	}

	span.LogFields(log.String("Organization", organizationAggregate.Organization.String()))
	return c.es.Save(ctx, organizationAggregate)
}
