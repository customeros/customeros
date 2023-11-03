package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LinkJobRoleCommandHandler interface {
	Handle(ctx context.Context, command *command.LinkJobRoleCommand) error
}

type linkJobRoleCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewLinkJobRoleCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) LinkJobRoleCommandHandler {
	return &linkJobRoleCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *linkJobRoleCommandHandler) Handle(ctx context.Context, cmd *command.LinkJobRoleCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "linkJobRoleCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	span.LogFields(log.Object("command", cmd))

	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	userAggregate, err := aggregate.LoadUserAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = userAggregate.LinkJobRole(ctx, cmd.Tenant, cmd.JobRoleId, cmd.LoggedInUserId); err != nil {
		return err
	}

	span.LogFields(log.String("User", userAggregate.User.String()))
	return c.es.Save(ctx, userAggregate)
}
