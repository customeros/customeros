package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type OrganizationTempAggregate struct {
	*aggregate.CommonTenantIdTempAggregate
}

func NewOrganizationTempAggregateWithTenantAndID(tenant, id string) *OrganizationTempAggregate {
	organizationTempAggregate := OrganizationTempAggregate{}
	organizationTempAggregate.CommonTenantIdTempAggregate = aggregate.NewCommonTempAggregateWithTenantAndId(OrganizationAggregateType, tenant, id)
	organizationTempAggregate.Tenant = tenant

	return &organizationTempAggregate
}

func (a *OrganizationTempAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationTempAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *organizationpb.EnrichOrganizationGrpcRequest:
		return nil, a.requestEnrichOrganization(ctx, r)
	case *organizationpb.RefreshDerivedDataGrpcRequest:
		return nil, a.refreshDerivedData(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *OrganizationTempAggregate) requestEnrichOrganization(ctx context.Context, request *organizationpb.EnrichOrganizationGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationTempAggregate.requestEnrichOrganization")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	enrichEvent, err := events.NewOrganizationRequestEnrich(a, request.Url)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRequestEnrich")
	}
	aggregate.EnrichEventWithMetadataExtended(&enrichEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(enrichEvent)
}

func (a *OrganizationTempAggregate) refreshDerivedData(ctx context.Context, request *organizationpb.RefreshDerivedDataGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationTempAggregate.refreshDerivedData")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	refreshDataEvent, err := events.NewOrganizationRefreshDerivedData(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRefreshDerivedData")
	}
	aggregate.EnrichEventWithMetadataExtended(&refreshDataEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(refreshDataEvent)
}
