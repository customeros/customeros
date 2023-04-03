package commands

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/phone_number"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type UpsertPhoneNumberCommandHandler interface {
	Handle(ctx context.Context, command *UpsertPhoneNumberCommand) error
}

type upsertPhoneNumberHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertPhoneNumberHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *upsertPhoneNumberHandler {
	return &upsertPhoneNumberHandler{log: log, cfg: cfg, es: es}
}

func (c *upsertPhoneNumberHandler) Handle(ctx context.Context, command *UpsertPhoneNumberCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertPhoneNumberHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", command.GetAggregateID()))

	if len(command.Tenant) == 0 {
		return phone_number.ErrMissingTenant
	}

	phoneNumberAggregate := aggregate.NewPhoneNumberAggregateWithTenantAndID(command.Tenant, command.AggregateID)
	err := c.es.Exists(ctx, phoneNumberAggregate.GetID())
	if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
		return err
	} else if err != nil && errors.Is(err, esdb.ErrStreamNotFound) {
		if err = phoneNumberAggregate.CreatePhoneNumber(ctx, command.Tenant, command.RawPhoneNumber, command.Source, command.SourceOfTruth, command.AppSource, command.CreatedAt, command.UpdatedAt); err != nil {
			return err
		}
	} else {
		if err = phoneNumberAggregate.UpdatePhoneNumber(ctx, command.Tenant, command.SourceOfTruth, command.UpdatedAt); err != nil {
			return err
		}
	}

	// FIXME alexb
	//span.LogFields(log.String("order", order.String()))
	return c.es.Save(ctx, phoneNumberAggregate)
}
