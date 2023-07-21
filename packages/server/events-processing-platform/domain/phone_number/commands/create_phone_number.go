package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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

func NewCreatePhoneNumberCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) CreatePhoneNumberCommandHandler {
	return &createPhoneNumberCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *createPhoneNumberCommandHandler) Handle(ctx context.Context, command *CreatePhoneNumberCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createPhoneNumberCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	phoneNumberAggregate, err := aggregate.LoadPhoneNumberAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = phoneNumberAggregate.CreatePhoneNumber(ctx, command.Tenant, command.RawPhoneNumber, command.Source.Source, command.Source.SourceOfTruth, command.Source.AppSource, command.CreatedAt, command.UpdatedAt); err != nil {
		return err
	}

	return h.es.Save(ctx, phoneNumberAggregate)
}
