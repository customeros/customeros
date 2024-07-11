package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/aggregate"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type UpsertLogEntryCommandHandler interface {
	Handle(ctx context.Context, command *cmd.UpsertLogEntryCommand) error
}

type upsertLogEntryCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewUpsertLogEntryCommandHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) UpsertLogEntryCommandHandler {
	return &upsertLogEntryCommandHandler{log: log, es: es, cfg: cfg}
}

func (c *upsertLogEntryCommandHandler) Handle(ctx context.Context, command *cmd.UpsertLogEntryCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "upsertLogEntryCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, command.Tenant, command.LoggedInUserId)
	span.LogFields(log.String("ObjectID", command.ObjectID))

	if err := validator.GetValidator().Struct(command); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for attempt := 0; attempt == 0 || attempt < c.cfg.RetriesOnOptimisticLockException; attempt++ {
		logEntryAggregate, err := aggregate.LoadLogEntryAggregate(ctx, c.es, command.Tenant, command.ObjectID)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if eventstore.IsAggregateNotFound(logEntryAggregate) {
			command.IsCreateCommand = true
		}
		if err = logEntryAggregate.HandleCommand(ctx, command); err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		err = c.es.Save(ctx, logEntryAggregate)
		if err == nil {
			return nil // Save successful
		}
		if eventstore.IsEventStoreErrorCodeWrongExpectedVersion(err) {
			// Handle concurrency error
			if attempt == c.cfg.RetriesOnOptimisticLockException-1 {
				// If we have reached the maximum number of retries, return an error
				tracing.TraceErr(span, err)
				return err
			}
			span.LogFields(log.Int("retryAttempt", attempt+1))
			time.Sleep(utils.BackOffExponentialDelay(attempt)) // backoffDelay is a function that increases the delay with each attempt
			continue                                           // Retry
		} else {
			// Some other error occurred
			tracing.TraceErr(span, err)
			return err
		}
	}
	err := errors.New("reached maximum number of retries")
	tracing.TraceErr(span, err)
	return err
}
