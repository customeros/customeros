package command_handler

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type UpdateOrganizationOwnerCommandHandler interface {
	Handle(ctx context.Context, command *command.UpdateOrganizationOwnerCommand) error
}

type updateOrganizationOwnerCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	ebs *eventbuffer.EventBufferStoreService
	cfg config.Utils
}

func NewUpdateOrganizationOwnerCommandHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils, ebs *eventbuffer.EventBufferStoreService) UpdateOrganizationOwnerCommandHandler {
	return &updateOrganizationOwnerCommandHandler{log: log, es: es, cfg: cfg, ebs: ebs}
}

func (h *updateOrganizationOwnerCommandHandler) Handle(ctx context.Context, cmd *command.UpdateOrganizationOwnerCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateOrganizationOwnerCommand.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	tracing.LogObjectAsJson(span, "command", cmd)

	validationError, done := validator.Validate(cmd, span)
	if done {
		return validationError
	}

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, h.es, cmd.Tenant, cmd.ObjectID, *eventstore.NewLoadAggregateOptions())
		if err != nil {
			return err
		}
		if err = organizationAggregate.HandleCommand(ctx, cmd); err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		// create notification event Park it in the event buffer for notifications
		event, err := createNotificationEvent(ctx, cmd, organizationAggregate)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		eventBufferUUID := fmt.Sprintf("%s-%s", cmd.ActorUserId, cmd.OrganizationId)
		err = h.ebs.Park(*event, cmd.Tenant, eventBufferUUID, time.Now().UTC().Add(time.Second*30))
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		// Save the aggregate to the event store
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

func createNotificationEvent(ctx context.Context, cmd *command.UpdateOrganizationOwnerCommand, aggregate *aggregate.OrganizationAggregate) (*eventstore.Event, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createNotificationEvent")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, cmd.Tenant, cmd.LoggedInUserId)
	tracing.LogObjectAsJson(span, "command", cmd)

	updatedAt := utils.Now()

	event, err := events.NewOrganizationOwnerUpdateNotificationEvent(aggregate, cmd.OwnerUserId, cmd.ActorUserId, cmd.OrganizationId, updatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "NewOrganizationOwnerUpdateNotificationEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: aggregate.GetTenant(),
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return &event, nil
}
