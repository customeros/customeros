package aggregate

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

// GetOrganizationAggregateID get organization aggregate id for eventstoredb
func GetOrganizationAggregateID(eventAggregateID string, tenant string) string {
	return strings.ReplaceAll(eventAggregateID, string(OrganizationAggregateType)+"-"+tenant+"-", "")
}

func IsAggregateNotFound(aggregate eventstore.Aggregate) bool {
	return aggregate.GetVersion() == 0
}

func LoadOrganizationAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, aggregateID string) (*OrganizationAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadOrganizationAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", aggregateID))

	organizationAggregate := NewOrganizationAggregateWithTenantAndID(tenant, aggregateID)

	err := eventStore.Exists(ctx, organizationAggregate.GetID())
	if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
		return nil, err
	}

	if err := eventStore.Load(ctx, organizationAggregate); err != nil {
		return nil, err
	}

	return organizationAggregate, nil
}
