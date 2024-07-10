package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func GetOrganizationObjectID(aggregateID string, tenant string) string {
	return aggregate.GetAggregateObjectID(aggregateID, tenant, OrganizationAggregateType)
}

func LoadOrganizationAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string, opts eventstore.LoadAggregateOptions) (*OrganizationAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadOrganizationAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	organizationAggregate := NewOrganizationAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, organizationAggregate, opts)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return organizationAggregate, nil
}

func LoadOrganizationTempAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string, opts eventstore.LoadAggregateOptions) (*OrganizationTempAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadOrganizationTempAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	organizationTempAggregate := NewOrganizationTempAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, organizationTempAggregate, opts)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return organizationTempAggregate, nil
}
