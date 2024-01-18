package event

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type CreateTenantBillingProfileEvent struct {
	Tenant       string             `json:"tenant" validate:"required"`
	Id           string             `json:"id" validate:"required"`
	CreatedAt    time.Time          `json:"createdAt"`
	SourceFields commonmodel.Source `json:"sourceFields"`
	Email        string             `json:"email"`
	Phone        string             `json:"phone"`
	AddressLine1 string             `json:"addressLine1"`
	AddressLine2 string             `json:"addressLine2"`
	AddressLine3 string             `json:"addressLine3"`
	LegalName    string             `json:"legalName"`
}

func NewCreateTenantBillingProfileEvent(aggregate eventstore.Aggregate, sourceFields commonmodel.Source, id, email, phone, addressLine1, addressLine2, addressLine3, legalName string, createdAt time.Time) (eventstore.Event, error) {
	eventData := CreateTenantBillingProfileEvent{
		Tenant:       aggregate.GetTenant(),
		Id:           id,
		CreatedAt:    createdAt,
		SourceFields: sourceFields,
		Email:        email,
		Phone:        phone,
		AddressLine1: addressLine1,
		AddressLine2: addressLine2,
		AddressLine3: addressLine3,
		LegalName:    legalName,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate CreateTenantBillingProfileEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, TenantAddBillingProfileV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for CreateTenantBillingProfileEvent")
	}

	return event, nil
}
