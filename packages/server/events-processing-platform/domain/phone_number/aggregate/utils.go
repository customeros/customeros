package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func GetPhoneNumberObjectID(aggregateID string, tenant string) string {
	return eventstore.GetAggregateObjectID(aggregateID, tenant, PhoneNumberAggregateType)
}

func LoadPhoneNumberAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*PhoneNumberAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadPhoneNumberAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	phoneNumberAggregate := NewPhoneNumberAggregateWithTenantAndID(tenant, objectID)

	err := eventstore.LoadAggregate(ctx, eventStore, phoneNumberAggregate, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return phoneNumberAggregate, nil
}
