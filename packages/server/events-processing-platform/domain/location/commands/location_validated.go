package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type LocationValidatedCommandHandler interface {
	Handle(ctx context.Context, command *LocationValidatedCommand) error
}

type locationValidatedCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewLocationValidatedCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *locationValidatedCommandHandler {
	return &locationValidatedCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *locationValidatedCommandHandler) Handle(ctx context.Context, command *LocationValidatedCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "locationValidatedCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	locationAggregate := aggregate.NewLocationAggregateWithTenantAndID(command.Tenant, command.ObjectID)
	err := c.es.Exists(ctx, locationAggregate.GetID())
	if err != nil && !errors.Is(err, eventstore.ErrAggregateNotFound) {
		return err
	}

	locationAddress := models.LocationAddress{}
	locationAddress.From(command.LocationAddressFields)

	locationAggregate, _ = aggregate.LoadLocationAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err = locationAggregate.LocationValidated(ctx, command.Tenant, command.RawAddress, command.CountryForValidation, locationAddress); err != nil {
		return err
	}
	return c.es.Save(ctx, locationAggregate)
}
