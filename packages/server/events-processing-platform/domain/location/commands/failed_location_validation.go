package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type FailedLocationValidationCommandHandler interface {
	Handle(ctx context.Context, command *FailedLocationValidationCommand) error
}

type failedLocationValidationCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewFailedLocationValidationCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *failedLocationValidationCommandHandler {
	return &failedLocationValidationCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *failedLocationValidationCommandHandler) Handle(ctx context.Context, command *FailedLocationValidationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "failedLocationValidationCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	locationAggregate := aggregate.NewLocationAggregateWithTenantAndID(command.Tenant, command.ObjectID)
	err := c.es.Exists(ctx, locationAggregate.GetID())

	if err != nil && !errors.Is(err, eventstore.ErrAggregateNotFound) {
		return err
	}

	locationAggregate, _ = aggregate.LoadLocationAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err = locationAggregate.FailLocationValidation(ctx, command.Tenant, command.RawAddress, command.Country, command.ValidationError); err != nil {
		return err
	}
	return c.es.Save(ctx, locationAggregate)
}
