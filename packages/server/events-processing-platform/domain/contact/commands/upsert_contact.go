package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpsertContactCommandHandler interface {
	Handle(ctx context.Context, command *UpsertContactCommand) error
}

type upsertContactCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertContactCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) UpsertContactCommandHandler {
	return &upsertContactCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *upsertContactCommandHandler) Handle(ctx context.Context, command *UpsertContactCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertContactCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ContactID", command.ObjectID))

	if command.Tenant == "" {
		tracing.TraceErr(span, eventstore.ErrMissingTenant)
		return eventstore.ErrMissingTenant
	}

	contactAggregate, err := aggregate.LoadContactAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if aggregate.IsAggregateNotFound(contactAggregate) {
		if err = contactAggregate.CreateContact(ctx, UpsertContactCommandToContactDto(command)); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		if err = contactAggregate.UpdateContact(ctx, UpsertContactCommandToContactDto(command)); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return h.es.Save(ctx, contactAggregate)
}
