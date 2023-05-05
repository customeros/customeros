package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type CreatePhoneNumberCommandHandler interface {
	Handle(ctx context.Context, command *CreatePhoneNumberCommand) error
}

type createPhoneNumberCommandHandler struct {
	log          logger.Logger
	cfg          *config.Config
	es           eventstore.AggregateStore
	repositories *repository.Repositories
}

func NewCreatePhoneNumberCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *createPhoneNumberCommandHandler {
	return &createPhoneNumberCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *createPhoneNumberCommandHandler) Handle(ctx context.Context, command *CreatePhoneNumberCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createPhoneNumberCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", command.GetAggregateID()))

	phoneNumberAggregate := aggregate.NewPhoneNumberAggregateWithTenantAndID(command.Tenant, command.AggregateID)
	err := h.es.Exists(ctx, phoneNumberAggregate.GetID())
	if err != nil && !errors.Is(err, eventstore.ErrAggregateNotFound) {
		return err
	}

	if err = phoneNumberAggregate.CreatePhoneNumber(ctx, command.Tenant, command.PhoneNumber, command.Source.Source, command.Source.SourceOfTruth, command.Source.AppSource, command.CreatedAt, command.UpdatedAt); err != nil {
		return err
	}

	span.LogFields(log.String("PhoneNumber", phoneNumberAggregate.PhoneNumber.String()))
	return h.es.Save(ctx, phoneNumberAggregate)
}
