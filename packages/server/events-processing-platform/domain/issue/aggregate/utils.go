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

func GetIssueObjectID(aggregateID string, tenant string) string {
	if tenant == "" {
		return getIssueObjectUUID(aggregateID)
	}
	return strings.ReplaceAll(aggregateID, string(IssueAggregateType)+"-"+tenant+"-", "")
}

// use this method when tenant is not known
func getIssueObjectUUID(aggregateID string) string {
	parts := strings.Split(aggregateID, "-")
	fullUUID := parts[len(parts)-5] + "-" + parts[len(parts)-4] + "-" + parts[len(parts)-3] + "-" + parts[len(parts)-2] + "-" + parts[len(parts)-1]
	return fullUUID
}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() < 0
}

func LoadIssueAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*IssueAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadIssueAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	issueAggregate := NewIssueAggregateWithTenantAndID(tenant, objectID)
	span.SetTag("aggregateID", issueAggregate.GetID())

	err := eventStore.Exists(ctx, issueAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			return nil, err
		} else {
			return issueAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, issueAggregate); err != nil {
		return nil, err
	}

	return issueAggregate, nil
}
