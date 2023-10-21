package command_handler

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type UpsertIssueCommandHandler interface {
	Handle(ctx context.Context, cmd *command.UpsertIssueCommand) error
}

type upsertIssueCommandHandler struct {
	log logger.Logger
	cfg *config.Config
	es  eventstore.AggregateStore
}

func NewUpsertIssueCommandHandler(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) UpsertIssueCommandHandler {
	return &upsertIssueCommandHandler{log: log, cfg: cfg, es: es}
}

func (c *upsertIssueCommandHandler) Handle(ctx context.Context, cmd *command.UpsertIssueCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertIssueCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.UserID)
	span.LogFields(log.String("command", fmt.Sprintf("%+v", cmd)))

	if err := validator.GetValidator().Struct(cmd); err != nil {
		wrappedErr := errors.Wrap(err, "failed validation for UpsertIssueCommand")
		tracing.TraceErr(span, wrappedErr)
		return wrappedErr
	}

	issueAggregate, err := aggregate.LoadIssueAggregate(ctx, c.es, cmd.Tenant, cmd.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if aggregate.IsAggregateNotFound(issueAggregate) {
		cmd.IsCreateCommand = true
	}
	if err = issueAggregate.HandleCommand(ctx, cmd); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return c.es.Save(ctx, issueAggregate)
}
