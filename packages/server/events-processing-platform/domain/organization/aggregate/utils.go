package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

func GetOrganizationObjectID(aggregateID string, tenant string) string {
	return strings.ReplaceAll(aggregateID, string(OrganizationAggregateType)+"-"+tenant+"-", "")
}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() <= 0
}

func LoadOrganizationAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*OrganizationAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadOrganizationAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("ObjectID", objectID))

	organizationAggregate := NewOrganizationAggregateWithTenantAndID(tenant, objectID)

	err := eventStore.Exists(ctx, organizationAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			return nil, err
		} else {
			return organizationAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, organizationAggregate); err != nil {
		return nil, err
	}

	return organizationAggregate, nil
}
