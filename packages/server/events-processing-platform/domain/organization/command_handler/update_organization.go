package command_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type UpdateOrganizationCommandHandler interface {
	Handle(ctx context.Context, command *command.UpdateOrganizationCommand) error
}

type updateOrganizationCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewUpdateOrganizationCommandHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) UpdateOrganizationCommandHandler {
	return &updateOrganizationCommandHandler{log: log, es: es, cfg: cfg}
}

func (h *updateOrganizationCommandHandler) Handle(ctx context.Context, cmd *command.UpdateOrganizationCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateOrganizationCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	tracing.LogObjectAsJson(span, "command", cmd)

	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	orgFields := &model.OrganizationFields{
		ID:                     cmd.ObjectID,
		Tenant:                 cmd.Tenant,
		OrganizationDataFields: cmd.DataFields,
		Source: cmnmod.Source{
			Source: cmd.Source,
		},
		UpdatedAt: cmd.UpdatedAt,
	}

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID)
		if err != nil {
			return err
		}
		if err = organizationAggregate.UpdateOrganization(ctx, orgFields, "", cmd.WebScrapedUrl, cmd.FieldsMask); err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		err = h.es.Save(ctx, organizationAggregate)
		if err == nil {
			return nil // Save successful
		}

		if eventstore.IsEventStoreErrorCodeWrongExpectedVersion(err) {
			// Handle concurrency error
			if attempt == h.cfg.RetriesOnOptimisticLockException-1 {
				// If we have reached the maximum number of retries, return an error
				tracing.TraceErr(span, err)
				return err
			}
			span.LogFields(log.Int("retryAttempt", attempt+1))
			time.Sleep(utils.BackOffExponentialDelay(attempt)) // backoffDelay is a function that increases the delay with each attempt
			continue                                           // Retry
		} else {
			// Some other error occurred
			tracing.TraceErr(span, errors.Wrap(err, "failed to save organization aggregate"))
			return err
		}
	}

	err := errors.New("reached maximum number of retries")
	tracing.TraceErr(span, err)
	return err
}
