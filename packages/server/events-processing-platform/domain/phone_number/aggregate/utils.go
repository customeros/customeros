package aggregate

import (
	es "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"strings"
)

// GetPhoneNumberAggregateID get phone_number aggregate id for eventstoredb
func GetPhoneNumberAggregateID(eventAggregateID string, tenant string) string {
	return strings.ReplaceAll(eventAggregateID, string(PhoneNumberAggregateType)+"-"+tenant+"-", "")
}

func IsAggregateNotFound(aggregate es.Aggregate) bool {
	return aggregate.GetVersion() == 0
}

func LoadPhoneNumberAggregate(ctx context.Context, eventStore es.AggregateStore, tenant, aggregateID string) (*PhoneNumberAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadPhoneNumberAggregate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", aggregateID))

	phoneNumberAggregate := NewPhoneNumberAggregateWithTenantAndID(tenant, aggregateID)

	err := eventStore.Exists(ctx, phoneNumberAggregate.GetID())
	if err != nil && !es.IsErrEsResourceNotFound(err) {
		return nil, err
	}

	if err := eventStore.Load(ctx, phoneNumberAggregate); err != nil {
		return nil, err
	}

	return phoneNumberAggregate, nil
}
