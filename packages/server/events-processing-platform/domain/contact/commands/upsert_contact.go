package commands

import (
	"context"
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

type upsertContactCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertContactCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *upsertContactCommandHandler {
	return &upsertContactCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *upsertContactCommandHandler) Handle(ctx context.Context, command *UpsertContactCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertContactCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", command.GetAggregateID()))

	if len(command.Tenant) == 0 {
		return eventstore.ErrMissingTenant
	}

	contactAggregate := aggregate.NewContactAggregateWithTenantAndID(command.Tenant, command.AggregateID)
	err := c.es.Exists(ctx, contactAggregate.GetID())
	if err != nil && !errors.Is(err, eventstore.ErrAggregateNotFound) {
		return err
	} else if err != nil && errors.Is(err, eventstore.ErrAggregateNotFound) {
		if err = contactAggregate.CreateContact(ctx, UpsertContactCommandToContactDto(command)); err != nil {
			return err
		}
	} else {
		contactAggregate, _ = aggregate.LoadContactAggregate(ctx, c.es, command.Tenant, command.AggregateID)
		if err = contactAggregate.UpdateContact(ctx, UpsertContactCommandToContactDto(command)); err != nil {
			return err
		}
	}

	span.LogFields(log.String("Contact", contactAggregate.Contact.String()))
	return c.es.Save(ctx, contactAggregate)
}
