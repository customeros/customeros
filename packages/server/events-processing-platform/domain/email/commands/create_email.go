package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type CreateEmailCommandHandler interface {
	Handle(ctx context.Context, command *CreateEmailCommand) error
}

type createEmailCommandHandler struct {
	log          logger.Logger
	cfg          *config.Config
	es           eventstore.AggregateStore
	repositories *repository.Repositories
}

func NewCreateEmailCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *createEmailCommandHandler {
	return &createEmailCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *createEmailCommandHandler) Handle(ctx context.Context, command *CreateEmailCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createEmailCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", command.GetAggregateID()))

	emailAggregate := aggregate.NewEmailAggregateWithTenantAndID(command.Tenant, command.AggregateID)
	err := c.es.Exists(ctx, emailAggregate.GetID())
	if err != nil && !errors.Is(err, eventstore.ErrAggregateNotFound) {
		return err
	}

	if err = emailAggregate.CreateEmail(ctx, command.Tenant, command.Email, command.Source.Source, command.Source.SourceOfTruth, command.Source.AppSource, command.CreatedAt, command.UpdatedAt); err != nil {
		return err
	}

	span.LogFields(log.String("Email", emailAggregate.Email.String()))
	return c.es.Save(ctx, emailAggregate)
}
