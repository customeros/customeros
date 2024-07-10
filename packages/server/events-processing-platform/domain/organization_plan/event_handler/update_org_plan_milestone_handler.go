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
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	orgplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/org_plan"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type UpdateOrganizationPlanMilestoneCommandHandler interface {
	Handle(ctx context.Context, baseRequest events.BaseRequest, request *orgplanpb.UpdateOrganizationPlanMilestoneGrpcRequest) error
}

type updateOrganizationPlanMilestoneCommandHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewUpdateOrganizationPlanMilestoneCommandHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) UpdateOrganizationPlanMilestoneCommandHandler {
	return &updateOrganizationPlanMilestoneCommandHandler{log: log, es: es, cfg: cfg}
}

// Handle processes the UpdateOrganizationPlanMilestoneCommand to update a new org plan.
func (h *updateOrganizationPlanMilestoneCommandHandler) Handle(ctx context.Context, baseRequest events.BaseRequest, request *orgplanpb.UpdateOrganizationPlanMilestoneGrpcRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateOrganizationPlanMilestoneCommandHandler.Handle")
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

		var dueDate time.Time
		if request.DueDate != nil {
			dueDate = request.DueDate.AsTime()
		}

		evt, err := event.NewOrganizationPlanMilestoneUpdateEvent(
			orgAggregate,
			request.OrganizationPlanId,
			baseRequest.ObjectID,
			request.Name,
			request.Order,
			GrpcItemsToDomainItems(request.Items),
			extractOrganizationPlanMilestoneFieldsMask(request.FieldsMask),
			request.Optional,
			request.Adhoc,
			request.Retired,
			updatedAtNotNil,
			dueDate,
			statusDetails,
		)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewOrganizationPlanMilestoneUpdateEvent")
		}

		commonAggregate.EnrichEventWithMetadataExtended(&evt, span, commonAggregate.EventMetadata{
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

func extractOrganizationPlanMilestoneFieldsMask(fields []orgplanpb.OrganizationPlanMilestoneFieldMask) []string {
	fieldsMask := make([]string, 0)
	if len(fields) == 0 {
		return fieldsMask
	}
	if containsOrganizationPlanMilestoneMaskFieldAll(fields) {
		return fieldsMask
	}
	for _, field := range fields {
		switch field {
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_NAME:
			fieldsMask = append(fieldsMask, event.FieldMaskName)
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_RETIRED:
			fieldsMask = append(fieldsMask, event.FieldMaskRetired)
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_ORDER:
			fieldsMask = append(fieldsMask, event.FieldMaskOrder)
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_OPTIONAL:
			fieldsMask = append(fieldsMask, event.FieldMaskOptional)
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_DUE_DATE:
			fieldsMask = append(fieldsMask, event.FieldMaskDueDate)
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_ITEMS:
			fieldsMask = append(fieldsMask, event.FieldMaskItems)
		case orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_STATUS_DETAILS:
			fieldsMask = append(fieldsMask, event.FieldMaskStatusDetails)
		}
	}
	return utils.RemoveDuplicates(fieldsMask)
}

func containsOrganizationPlanMilestoneMaskFieldAll(fields []orgplanpb.OrganizationPlanMilestoneFieldMask) bool {
	for _, field := range fields {
		if field == orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_ALL {
			return true
		}
	}
	return false
}

func GrpcItemsToDomainItems(items []*orgplanpb.OrganizationPlanMilestoneItem) []model.OrganizationPlanMilestoneItem {
	domainItems := make([]model.OrganizationPlanMilestoneItem, 0)
	for _, item := range items {
		domainItems = append(domainItems, model.OrganizationPlanMilestoneItem{
			Status:    item.Status,
			Text:      item.Text,
			UpdatedAt: item.UpdatedAt.AsTime(),
			Uuid:      item.Uuid,
		})
	}
	return domainItems
}
