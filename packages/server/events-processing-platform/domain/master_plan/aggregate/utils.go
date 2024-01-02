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

func GetMasterPlanObjectID(aggregateID string, tenant string) string {
	if tenant == "" {
		return getMasterPlanObjectUUID(aggregateID)
	}
	return strings.ReplaceAll(aggregateID, string(MasterPlanAggregateType)+"-"+tenant+"-", "")
}

// getMasterPlanObjectUUID generates the UUID for a master plan when the tenant is not known.
func getMasterPlanObjectUUID(aggregateID string) string {
	parts := strings.Split(aggregateID, "-")
	fullUUID := parts[len(parts)-5] + "-" + parts[len(parts)-4] + "-" + parts[len(parts)-3] + "-" + parts[len(parts)-2] + "-" + parts[len(parts)-1]
	return fullUUID
}

// IsAggregateNotFound checks if the provided aggregate is not found.
func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() < 0
}

func LoadMasterPlanAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*MasterPlanAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadMasterPlanAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	masterPlanAggregate := NewMasterPlanAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, masterPlanAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			tracing.TraceErr(span, err)
			return nil, err
		} else {
			return masterPlanAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, masterPlanAggregate); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return masterPlanAggregate, nil
}
