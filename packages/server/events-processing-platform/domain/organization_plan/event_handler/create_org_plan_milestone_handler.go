package event_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	commonAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
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

type CreateOrganizationPlanMilestoneHandler interface {
	Handle(ctx context.Context, baseRequest events.BaseRequest, request *orgplanpb.CreateOrganizationPlanMilestoneGrpcRequest) error
}

type createOrganizationPlanMilestoneHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewCreateOrganizationPlanMilestoneHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) CreateOrganizationPlanMilestoneHandler {
	return &createOrganizationPlanMilestoneHandler{log: log, es: es, cfg: cfg}
}

// Handle processes the CreateOrganizationPlanMilestone event to create a new org plan.
func (h *createOrganizationPlanMilestoneHandler) Handle(ctx context.Context, baseRequest events.BaseRequest, request *orgplanpb.CreateOrganizationPlanMilestoneGrpcRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createOrganizationPlanMilestoneHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, baseRequest.Tenant, baseRequest.LoggedInUserId)
	tracing.LogObjectAsJson(span, "common", baseRequest)
	tracing.LogObjectAsJson(span, "request", request)

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		// Load or initialize the org aggregate
		organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, h.es, baseRequest.Tenant, request.OrgId, eventstore.LoadAggregateOptions{})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if eventstore.IsAggregateNotFound(organizationAggregate) {
			tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
			return eventstore.ErrAggregateNotFound
		}

		createdAtNotNil := utils.IfNotNilTimeWithDefault(request.CreatedAt, utils.Now())

		evt, err := event.NewOrganizationPlanMilestoneCreateEvent(organizationAggregate, request.OrganizationPlanId, baseRequest.ObjectID, request.Name, request.Order, request.Items, request.Optional, request.Adhoc, baseRequest.SourceFields, createdAtNotNil, request.DueDate.AsTime())
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewOrganizationPlanMilestoneCreateEvent")
		}

		commonAggregate.EnrichEventWithMetadataExtended(&evt, span, commonAggregate.EventMetadata{
			Tenant: request.Tenant,
			UserId: request.LoggedInUserId,
			App:    request.SourceFields.AppSource,
		})

		err = organizationAggregate.Apply(evt)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		// Persist the changes to the event store
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
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
