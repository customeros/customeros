package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpsertLocationCommandHandler interface {
	Handle(ctx context.Context, command *UpsertLocationCommand) error
}

type upsertLocationHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertLocationHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) UpsertLocationCommandHandler {
	return &upsertLocationHandler{log: log, cfg: cfg, es: es}
}

func (h *upsertLocationHandler) Handle(ctx context.Context, command *UpsertLocationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertLocationHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if len(command.Tenant) == 0 {
		return eventstore.ErrMissingTenant
	}

	locationAggregate, err := aggregate.LoadLocationAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	locationAddress := models.LocationAddress{}
	locationAddress.From(command.LocationAddressFields)
	if aggregate.IsAggregateNotFound(locationAggregate) {
		if err = locationAggregate.CreateLocation(ctx, command.Tenant, command.Name, command.RawAddress, locationAddress, command.Source, command.CreatedAt, command.UpdatedAt); err != nil {
			return err
		}
	} else {
		if err = locationAggregate.UpdateLocation(ctx, command.Tenant, command.Name, command.RawAddress, locationAddress, command.Source.SourceOfTruth, command.UpdatedAt); err != nil {
			return err
		}
	}

	return h.es.Save(ctx, locationAggregate)
}
