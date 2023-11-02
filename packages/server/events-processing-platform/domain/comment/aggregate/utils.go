package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

func GetCommentObjectID(aggregateID string, tenant string) string {
	if tenant == "" {
		return getCommentObjectUUID(aggregateID)
	}
	return strings.ReplaceAll(aggregateID, string(CommentAggregateType)+"-"+tenant+"-", "")
}

// use this method when tenant is not known
func getCommentObjectUUID(aggregateID string) string {
	parts := strings.Split(aggregateID, "-")
	fullUUID := parts[len(parts)-5] + "-" + parts[len(parts)-4] + "-" + parts[len(parts)-3] + "-" + parts[len(parts)-2] + "-" + parts[len(parts)-1]
	return fullUUID
}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() < 0
}

func LoadCommentAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*CommentAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadCommentAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	commentAggregate := NewCommentAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, commentAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			tracing.TraceErr(span, err)
			return nil, err
		} else {
			return commentAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, commentAggregate); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return commentAggregate, nil
}
