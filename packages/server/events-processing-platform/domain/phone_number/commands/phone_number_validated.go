package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type PhoneNumberValidatedCommandHandler interface {
	Handle(ctx context.Context, command *PhoneNumberValidatedCommand) error
}

type phoneNumberValidatedCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewPhoneNumberValidatedCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *phoneNumberValidatedCommandHandler {
	return &phoneNumberValidatedCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *phoneNumberValidatedCommandHandler) Handle(ctx context.Context, command *PhoneNumberValidatedCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "phoneNumberValidatedCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", command.GetAggregateID()))

	phoneNumberAggregate := aggregate.NewPhoneNumberAggregateWithTenantAndID(command.Tenant, command.AggregateID)
	err := c.es.Exists(ctx, phoneNumberAggregate.GetID())
	if err != nil && !errors.Is(err, eventstore.ErrAggregateNotFound) {
		return err
	}

	phoneNumberAggregate, _ = aggregate.LoadPhoneNumberAggregate(ctx, c.es, command.Tenant, command.AggregateID)
	if err = phoneNumberAggregate.PhoneNumberValidated(ctx, command.Tenant, command.RawPhoneNumber, command.E164, command.CountryCodeA2); err != nil {
		return err
	}
	return c.es.Save(ctx, phoneNumberAggregate)
}
