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

type SkippedLocationValidationCommandHandler interface {
	Handle(ctx context.Context, command *SkippedLocationValidationCommand) error
}

type skippedLocationValidationCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewSkippedLocationValidationCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *skippedLocationValidationCommandHandler {
	return &skippedLocationValidationCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *skippedLocationValidationCommandHandler) Handle(ctx context.Context, command *SkippedLocationValidationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "skippedLocationValidationCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	locationAggregate := aggregate.NewLocationAggregateWithTenantAndID(command.Tenant, command.ObjectID)
	err := c.es.Exists(ctx, locationAggregate.GetID())

	if err != nil && !errors.Is(err, eventstore.ErrAggregateNotFound) {
		return err
	}

	locationAggregate, _ = aggregate.LoadLocationAggregate(ctx, c.es, command.Tenant, command.ObjectID)
	if err = locationAggregate.SkipLocationValidation(ctx, command.Tenant, command.RawAddress, command.ValidationSkipReason); err != nil {
		return err
	}
	return c.es.Save(ctx, locationAggregate)
}
