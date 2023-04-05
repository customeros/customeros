package commands

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type UpsertContactCommandHandler interface {
	Handle(ctx context.Context, command *UpsertContactCommand) error
}

type upsertContactHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertContactHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *upsertContactHandler {
	return &upsertContactHandler{log: log, cfg: cfg, es: es}
}

func (c *upsertContactHandler) Handle(ctx context.Context, command *UpsertContactCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertContactHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", command.GetAggregateID()))

	if len(command.Tenant) == 0 {
		return eventstore.ErrMissingTenant
	}

	contactAggregate := aggregate.NewContactAggregateWithTenantAndID(command.Tenant, command.AggregateID)
	err := c.es.Exists(ctx, contactAggregate.GetID())
	if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
		return err
	} else if err != nil && errors.Is(err, esdb.ErrStreamNotFound) {
		if err = contactAggregate.CreateContact(ctx, command.Tenant, command.FirstName, command.LastName, command.Name, command.Prefix, command.Source.Source, command.Source.SourceOfTruth, command.Source.AppSource, command.CreatedAt, command.UpdatedAt); err != nil {
			return err
		}
	} else {
		contactAggregate, _ = aggregate.LoadContactAggregate(ctx, c.es, command.Tenant, command.AggregateID)
		if err = contactAggregate.UpdateContact(ctx, command.Tenant, command.FirstName, command.LastName, command.Name, command.Prefix, command.Source.SourceOfTruth, command.UpdatedAt); err != nil {
			return err
		}
	}

	span.LogFields(log.String("Contact", contactAggregate.Contact.String()))
	return c.es.Save(ctx, contactAggregate)
}
