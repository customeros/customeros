package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

func GetLogEntryObjectID(aggregateID string, tenant string) string {
	if tenant == "" {
		return getLogEntryObjectUUID(aggregateID)
	}
	return strings.ReplaceAll(aggregateID, string(LogEntryAggregateType)+"-"+tenant+"-", "")
}

// use this method when tenant is not known
func getLogEntryObjectUUID(aggregateID string) string {
	parts := strings.Split(aggregateID, "-")
	fullUUID := parts[len(parts)-5] + "-" + parts[len(parts)-4] + "-" + parts[len(parts)-3] + "-" + parts[len(parts)-2] + "-" + parts[len(parts)-1]
	return fullUUID
}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() < 0
}

func LoadLogEntryAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*LogEntryAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadLogEntryAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("ObjectID", objectID))

	logEntryAggregate := NewLogEntryAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, logEntryAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			return nil, err
		} else {
			return logEntryAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, logEntryAggregate); err != nil {
		return nil, err
	}

	return logEntryAggregate, nil
}
