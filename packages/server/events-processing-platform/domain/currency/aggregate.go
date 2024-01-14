package currency

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
	CurrencyAggregateType eventstore.AggregateType = "currency"
)

type currencyAggregate struct {
	*aggregate.CommonIdAggregate
	Currency *Currency
}

func GetCurrencyObjectID(aggregateID string) string {
	return getCurrencyObjectUUID(aggregateID)
}

func getCurrencyObjectUUID(aggregateID string) string {
	parts := strings.Split(aggregateID, "-")
	fullUUID := parts[len(parts)-5] + "-" + parts[len(parts)-4] + "-" + parts[len(parts)-3] + "-" + parts[len(parts)-2] + "-" + parts[len(parts)-1]
	return fullUUID
}

func LoadCurrencyAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*currencyAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadCurrencyAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	currencyAggregate := NewCurrencyAggregateWithID(objectID)

	err := eventStore.Exists(ctx, currencyAggregate.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			tracing.TraceErr(span, err)
			return nil, err
		} else {
			return currencyAggregate, nil
		}
	}

	if err = eventStore.Load(ctx, currencyAggregate); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return currencyAggregate, nil
}

func NewCurrencyAggregateWithID(id string) *currencyAggregate {
	currencyAggregate := currencyAggregate{}
	currencyAggregate.CommonIdAggregate = aggregate.NewCommonAggregateWithId(CurrencyAggregateType, id)
	currencyAggregate.SetWhen(currencyAggregate.When)
	currencyAggregate.Currency = &Currency{}

	return &currencyAggregate
}

func (a *currencyAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case CurrencyCreateV1:
		return a.onCurrencyCreate(evt)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *currencyAggregate) onCurrencyCreate(evt eventstore.Event) error {
	var eventData CurrencyCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Currency.CreatedAt = eventData.CreatedAt
	a.Currency.Name = eventData.Name
	a.Currency.Symbol = eventData.Symbol

	return nil
}
