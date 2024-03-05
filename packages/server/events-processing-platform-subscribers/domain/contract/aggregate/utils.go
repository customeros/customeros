package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func GetContractObjectID(aggregateID, tenant string) string {
	return aggregate.GetAggregateObjectID(aggregateID, tenant, ContractAggregateType)
}

func LoadContractAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string, options eventstore.LoadAggregateOptions) (*ContractAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadContractAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	contractAggregate := NewContractAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, contractAggregate, options)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return contractAggregate, nil
}

func LoadContractTempAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*ContractTempAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadContractTempAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	contractTempAggregate := NewContractTempAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, contractTempAggregate, eventstore.LoadAggregateOptions{SkipLoadEvents: true})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return contractTempAggregate, nil
}
