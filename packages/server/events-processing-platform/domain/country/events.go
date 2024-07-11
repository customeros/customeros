package country

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	CountryCreateV1 = "V1_COUNTRY_CREATE"
)

type CountryCreateEvent struct {
	CreatedAt time.Time `json:"createdAt"`

	Name      string `json:"name" validate:"required"`
	CodeA2    string `json:"codeA2" validate:"required"`
	CodeA3    string `json:"codeA3" validate:"required"`
	PhoneCode string `json:"phoneCode" validate:"required"`
}

func NewCountryCreateEvent(aggregate eventstore.Aggregate, name, codeA2, codeA3, phoneCode string, createdAt time.Time) (eventstore.Event, error) {
	eventData := CountryCreateEvent{
		CreatedAt: createdAt,

		Name:      name,
		CodeA2:    codeA2,
		CodeA3:    codeA3,
		PhoneCode: phoneCode,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate CountryCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, CountryCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for CountryCreateEvent")
	}

	return event, nil
}
