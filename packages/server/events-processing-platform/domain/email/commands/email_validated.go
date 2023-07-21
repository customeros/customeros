package commands

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type EmailValidatedCommandHandler interface {
	Handle(ctx context.Context, command *EmailValidatedCommand) error
}

type emailValidatedCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewEmailValidatedCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) EmailValidatedCommandHandler {
	return &emailValidatedCommandHandler{log: log, cfg: cfg, es: es}
}

func (h *emailValidatedCommandHandler) Handle(ctx context.Context, command *EmailValidatedCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailValidatedCommandHandler.Handle")
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.Tenant), log.String("ObjectID", command.ObjectID))

	emailAggregate, err := aggregate.LoadEmailAggregate(ctx, h.es, command.Tenant, command.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = emailAggregate.EmailValidated(ctx, command.Tenant, command.RawEmail, command.IsReachable, command.ValidationError, command.Domain, command.Username, command.EmailAddress,
		command.AcceptsMail, command.CanConnectSmtp, command.HasFullInbox, command.IsCatchAll, command.IsDeliverable, command.IsDisabled, command.IsValidSyntax); err != nil {
		return err
	}
	return h.es.Save(ctx, emailAggregate)
}
