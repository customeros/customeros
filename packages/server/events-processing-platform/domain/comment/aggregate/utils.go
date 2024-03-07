package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func GetCommentObjectID(aggregateID string, tenant string) string {
	return aggregate.GetAggregateObjectID(aggregateID, tenant, CommentAggregateType)
}

func LoadCommentAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*CommentAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadCommentAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	commentAggregate := NewCommentAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, commentAggregate, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return commentAggregate, nil
}
