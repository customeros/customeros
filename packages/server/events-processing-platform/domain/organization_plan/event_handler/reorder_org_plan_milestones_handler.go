package event_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	event "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	orgplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/org_plan"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ReorderOrganizationPlanMilestonesHandler interface {
	Handle(ctx context.Context, baseRequest events.BaseRequest, request *orgplanpb.ReorderOrganizationPlanMilestonesGrpcRequest) error
}

type reorderOrganizationPlanMilestonesHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewReorderOrganizationPlanMilestonesHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) ReorderOrganizationPlanMilestonesHandler {
	return &reorderOrganizationPlanMilestonesHandler{log: log, es: es, cfg: cfg}
}

func (h *reorderOrganizationPlanMilestonesHandler) Handle(ctx context.Context, baseRequest events.BaseRequest, request *orgplanpb.ReorderOrganizationPlanMilestonesGrpcRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReorderOrganizationPlanMilestonesHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, baseRequest.Tenant, baseRequest.LoggedInUserId)
	tracing.LogObjectAsJson(span, "common", baseRequest)
	tracing.LogObjectAsJson(span, "request", request)

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		// Load or initialize the org aggregate
		orgAggregate, err := aggregate.LoadOrganizationAggregate(ctx, h.es, baseRequest.Tenant, request.OrgId, eventstore.LoadAggregateOptions{})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if eventstore.IsAggregateNotFound(orgAggregate) {
			tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
			return eventstore.ErrAggregateNotFound
		}

		updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), utils.Now())

		evt, err := event.NewOrganizationPlanMilestoneReorderEvent(orgAggregate, request.OrganizationPlanId, request.OrganizationPlanMilestoneIds, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewOrganizationPlanMilestoneReorderEvent")
		}

		eventstore.EnrichEventWithMetadataExtended(&evt, span, eventstore.EventMetadata{
			Tenant: request.Tenant,
			UserId: request.LoggedInUserId,
			App:    baseRequest.SourceFields.AppSource,
		})

		err = orgAggregate.Apply(evt)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		// Persist the changes to the event store
		err = h.es.Save(ctx, orgAggregate)
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
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
