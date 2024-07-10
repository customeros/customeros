package reminder

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func GetReminderObjectID(aggregateID string, tenant string) string {
	return aggregate.GetAggregateObjectID(aggregateID, tenant, ReminderAggregateType)
}

func LoadReminderAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string, opts eventstore.LoadAggregateOptions) (*ReminderAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadReminderAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	reminderAggregate := NewReminderAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, reminderAggregate, opts)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return reminderAggregate, nil
}
