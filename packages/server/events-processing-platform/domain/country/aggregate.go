package country

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

const (
	CountryAggregateType eventstore.AggregateType = "country"
)

type countryAggregate struct {
	*aggregate.CommonIdAggregate
	Country *Country
}

func GetCountryObjectID(aggregateID string) string {
	return getCountryObjectUUID(aggregateID)
}

func getCountryObjectUUID(aggregateID string) string {
	parts := strings.Split(aggregateID, "-")
	fullUUID := parts[len(parts)-5] + "-" + parts[len(parts)-4] + "-" + parts[len(parts)-3] + "-" + parts[len(parts)-2] + "-" + parts[len(parts)-1]
	return fullUUID
}

func LoadCountryAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*countryAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadCountryAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	countryAggregate := NewCountryAggregateWithID(objectID)

	err := eventStore.Exists(ctx, countryAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			tracing.TraceErr(span, err)
			return nil, err
		} else {
			return countryAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, countryAggregate); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return countryAggregate, nil
}

func NewCountryAggregateWithID(id string) *countryAggregate {
	countryAggregate := countryAggregate{}
	countryAggregate.CommonIdAggregate = aggregate.NewCommonAggregateWithId(CountryAggregateType, id)
	countryAggregate.SetWhen(countryAggregate.When)
	countryAggregate.Country = &Country{}

	return &countryAggregate
}

func (a *countryAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case CountryCreateV1:
		return a.onCountryCreate(evt)
	default:
		if strings.HasPrefix(evt.GetEventType(), "$") {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *countryAggregate) onCountryCreate(evt eventstore.Event) error {
	var eventData CountryCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Country.CreatedAt = eventData.CreatedAt
	a.Country.Name = eventData.Name
	a.Country.CodeA2 = eventData.CodeA2
	a.Country.CodeA3 = eventData.CodeA3
	a.Country.PhoneCode = eventData.PhoneCode

	return nil
}
