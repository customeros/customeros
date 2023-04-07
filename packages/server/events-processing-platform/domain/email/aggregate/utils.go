package aggregate

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	es "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"strings"
)

// GetEmailAggregateID get email aggregate id for eventstoredb
func GetEmailAggregateID(eventAggregateID string, tenant string) string {
	return strings.ReplaceAll(eventAggregateID, string(EmailAggregateType)+"-"+tenant+"-", "")
}

func IsAggregateNotFound(aggregate es.Aggregate) bool {
	return aggregate.GetVersion() == 0
}

func LoadEmailAggregate(ctx context.Context, eventStore es.AggregateStore, tenant, aggregateID string) (*EmailAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadEmailAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", aggregateID))

	emailAggregate := NewEmailAggregateWithTenantAndID(tenant, aggregateID)

	err := eventStore.Exists(ctx, emailAggregate.GetID())
	if err != nil && !errors.Is(err, esdb.ErrStreamNotFound) {
		return nil, err
	}

	if err := eventStore.Load(ctx, emailAggregate); err != nil {
		return nil, err
	}

	return emailAggregate, nil
}
