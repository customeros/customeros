package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type UpsertOrganizationCommandHandler interface {
	Handle(ctx context.Context, command *UpsertOrganizationCommand) error
}

type upsertOrganizationCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertOrganizationCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *upsertOrganizationCommandHandler {
	return &upsertOrganizationCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *upsertOrganizationCommandHandler) Handle(ctx context.Context, command *UpsertOrganizationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertOrganizationCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if len(command.Tenant) == 0 {
		return eventstore.ErrMissingTenant
	}

	organizationAggregate := aggregate.NewOrganizationAggregateWithTenantAndID(command.Tenant, command.ObjectID)
	err := c.es.Exists(ctx, organizationAggregate.GetID())
	if err != nil && !errors.Is(err, eventstore.ErrAggregateNotFound) {
		return err
	} else if err != nil && errors.Is(err, eventstore.ErrAggregateNotFound) {
		if err = organizationAggregate.CreateOrganization(ctx, UpsertOrganizationCommandToOrganizationDto(command)); err != nil {
			return err
		}
	} else {
		organizationAggregate, _ = aggregate.LoadOrganizationAggregate(ctx, c.es, command.Tenant, command.ObjectID)
		if err = organizationAggregate.UpdateOrganization(ctx, UpsertOrganizationCommandToOrganizationDto(command)); err != nil {
			return err
		}
	}

	span.LogFields(log.String("Organization", organizationAggregate.Organization.String()))
	return c.es.Save(ctx, organizationAggregate)
}
