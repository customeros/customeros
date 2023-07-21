package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UpsertPhoneNumberCommandHandler interface {
	Handle(ctx context.Context, command *UpsertPhoneNumberCommand) error
}

type upsertPhoneNumberHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertPhoneNumberHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) UpsertPhoneNumberCommandHandler {
	return &upsertPhoneNumberHandler{log: log, cfg: cfg, es: es}
}

func (h *upsertPhoneNumberHandler) Handle(ctx context.Context, command *UpsertPhoneNumberCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpsertPhoneNumberHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	if command.Tenant == "" {
		return eventstore.ErrMissingTenant
	}

	phoneNumberAggregate, err := aggregate.LoadPhoneNumberAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if aggregate.IsAggregateNotFound(phoneNumberAggregate) {
		if err = phoneNumberAggregate.CreatePhoneNumber(ctx, command.Tenant, command.RawPhoneNumber, command.Source.Source, command.Source.SourceOfTruth, command.Source.AppSource, command.CreatedAt, command.UpdatedAt); err != nil {
			return err
		}
	} else {
		if err = phoneNumberAggregate.UpdatePhoneNumber(ctx, command.Tenant, command.Source.SourceOfTruth, command.UpdatedAt); err != nil {
			return err
		}
	}

	return h.es.Save(ctx, phoneNumberAggregate)
}
