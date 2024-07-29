package event_handler

import (
	"context"
	events "github.com/openline-ai/openline-customer-os/packages/server/events/event"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	orgplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/org_plan"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type UpdateOrganizationPlanHandler interface {
	Handle(ctx context.Context, baseRequest events.BaseRequest, request *orgplanpb.UpdateOrganizationPlanGrpcRequest) error
}

type updateOrganizationPlanHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewUpdateOrganizationPlanHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) UpdateOrganizationPlanHandler {
	return &updateOrganizationPlanHandler{log: log, es: es, cfg: cfg}
}

// Handle processes the UpdateOrganizationPlanCommand to update a new master plan.
func (h *updateOrganizationPlanHandler) Handle(ctx context.Context, baseRequest events.BaseRequest, request *orgplanpb.UpdateOrganizationPlanGrpcRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateOrganizationPlanHandler.Handle")
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
		statusDetails := model.OrganizationPlanDetails{}
		if request.StatusDetails != nil {
			statusDetails = model.OrganizationPlanDetails{
				Status:    request.StatusDetails.Status,
				UpdatedAt: updatedAtNotNil,
				Comments:  request.StatusDetails.Comments,
			}
		}

		evt, err := event.NewOrganizationPlanUpdateEvent(orgAggregate, request.OrganizationPlanId, request.Name, request.Retired, updatedAtNotNil, extractOrganizationPlanFieldsMask(request.FieldsMask), statusDetails)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewOrganizationPlanUpdateEvent")
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

func extractOrganizationPlanFieldsMask(fields []orgplanpb.OrganizationPlanFieldMask) []string {
	fieldsMask := make([]string, 0)
	if len(fields) == 0 {
		return fieldsMask
	}
	if containsOrganizationPlanMaskFieldAll(fields) {
		return fieldsMask
	}
	for _, field := range fields {
		switch field {
		case orgplanpb.OrganizationPlanFieldMask_ORGANIZATION_PLAN_PROPERTY_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskName)
		case orgplanpb.OrganizationPlanFieldMask_ORGANIZATION_PLAN_PROPERTY_RETIRED:
			fieldsMask = append(fieldsMask, event.FieldMaskRetired)
		case orgplanpb.OrganizationPlanFieldMask_ORGANIZATION_PLAN_PROPERTY_STATUS_DETAILS:
			fieldsMask = append(fieldsMask, event.FieldMaskStatusDetails)
		}
	}
	return utils.RemoveDuplicates(fieldsMask)
}

func containsOrganizationPlanMaskFieldAll(fields []orgplanpb.OrganizationPlanFieldMask) bool {
	for _, field := range fields {
		if field == orgplanpb.OrganizationPlanFieldMask_ORGANIZATION_PLAN_PROPERTY_ALL {
			return true
		}
	}
	return false
}
