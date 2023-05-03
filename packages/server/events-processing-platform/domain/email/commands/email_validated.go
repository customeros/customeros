package commands

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type EmailValidatedCommandHandler interface {
	Handle(ctx context.Context, command *EmailValidatedCommand) error
}

type emailValidatedCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewEmailValidatedCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *emailValidatedCommandHandler {
	return &emailValidatedCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *emailValidatedCommandHandler) Handle(ctx context.Context, command *EmailValidatedCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "emailValidatedCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("AggregateID", command.GetAggregateID()))

	emailAggregate := aggregate.NewEmailAggregateWithTenantAndID(command.Tenant, command.AggregateID)
	err := c.es.Exists(ctx, emailAggregate.GetID())
	if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
		return err
	}

	if err = emailAggregate.EmailValidated(ctx, command.Tenant, command.Email, command.ValidationError, command.Domain, command.Username, command.NormalizedEmail,
		command.AcceptsMail, command.CanConnectSmtp, command.HasFullInbox, command.IsCatchAll, command.IsDeliverable, command.IsDisabled, command.IsValidSyntax); err != nil {
		return err
	}
	return c.es.Save(ctx, emailAggregate)
}
