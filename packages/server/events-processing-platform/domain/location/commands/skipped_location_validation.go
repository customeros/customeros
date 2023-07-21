package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SkippedLocationValidationCommandHandler interface {
	Handle(ctx context.Context, command *SkippedLocationValidationCommand) error
}

type skippedLocationValidationCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewSkippedLocationValidationCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) SkippedLocationValidationCommandHandler {
	return &skippedLocationValidationCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *skippedLocationValidationCommandHandler) Handle(ctx context.Context, command *SkippedLocationValidationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "skippedLocationValidationCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	locationAggregate, err := aggregate.LoadLocationAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = locationAggregate.SkipLocationValidation(ctx, command.Tenant, command.RawAddress, command.ValidationSkipReason); err != nil {
		return err
	}
	return h.es.Save(ctx, locationAggregate)
}
