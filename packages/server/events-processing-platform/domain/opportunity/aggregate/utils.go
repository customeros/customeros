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

func GetOpportunityObjectID(aggregateID string, tenant string) string {
	if tenant == "" {
		return getOpportunityObjectUUID(aggregateID)
	}
	return strings.ReplaceAll(aggregateID, string(OpportunityAggregateType)+"-"+tenant+"-", "")
}

// use this method when tenant is not known
func getOpportunityObjectUUID(aggregateID string) string {
	parts := strings.Split(aggregateID, "-")
	fullUUID := parts[len(parts)-5] + "-" + parts[len(parts)-4] + "-" + parts[len(parts)-3] + "-" + parts[len(parts)-2] + "-" + parts[len(parts)-1]
	return fullUUID
}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() < 0
}

func LoadOpportunityAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*OpportunityAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadOpportunityAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	opportunityAggregate := NewOpportunityAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, opportunityAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			tracing.TraceErr(span, err)
			return nil, err
		} else {
			return opportunityAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, opportunityAggregate); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return opportunityAggregate, nil
}
