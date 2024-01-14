package currency

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

const (
	CurrencyCreateV1 = "V1_CURRENCY_CREATE"
)

type CurrencyCreateEvent struct {
	CreatedAt    time.Time          `json:"createdAt"`
	SourceFields commonmodel.Source `json:"sourceFields"`

	Name   string `json:"name" validate:"required"`
	Symbol string `json:"symbol" validate:"required"`
}

func NewCurrencyCreateEvent(aggregate eventstore.Aggregate, name, symbol string, createdAt time.Time, sourceFields commonmodel.Source) (eventstore.Event, error) {
	eventData := CurrencyCreateEvent{
		CreatedAt:    createdAt,
		SourceFields: sourceFields,

		Name:   name,
		Symbol: symbol,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate CurrencyCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, CurrencyCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for CurrencyCreateEvent")
	}

	return event, nil
}

type EventHandlers struct {
	CurrencyCreate CurrencyCreateHandler
}

func NewEventHandlers(log logger.Logger, es eventstore.AggregateStore) *EventHandlers {
	return &EventHandlers{
		CurrencyCreate: NewCurrencyCreateHandler(log, es),
	}
}
