package event_handler

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	event "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	orgplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/org_plan"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type CreateOrganizationPlanHandler interface {
	Handle(ctx context.Context, baseRequest events.BaseRequest, request *orgplanpb.CreateOrganizationPlanGrpcRequest) error
}

type createOrganizationPlanHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewCreateOrganizationPlanHandler(log logger.Logger, es eventstore.AggregateStore) CreateOrganizationPlanHandler {
	return &createOrganizationPlanHandler{log: log, es: es}
}

// Handle processes the CreateOrganizationPlanCommand to create a new master plan.
func (h *createOrganizationPlanHandler) Handle(ctx context.Context, baseRequest events.BaseRequest, request *orgplanpb.CreateOrganizationPlanGrpcRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createOrganizationPlanCommandHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, baseRequest.Tenant, baseRequest.LoggedInUserId)
	tracing.LogObjectAsJson(span, "common", baseRequest)
	tracing.LogObjectAsJson(span, "request", request)

	// Load or initialize the org plan aggregate
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

	createEvent, err := event.NewOrganizationPlanCreateEvent(organizationAggregate, baseRequest.ObjectID, request.MasterPlanId, request.OrgId, request.Name, baseRequest.SourceFields, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationPlanCreateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&createEvent, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	err = organizationAggregate.Apply(createEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err = h.es.Save(ctx, organizationAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
