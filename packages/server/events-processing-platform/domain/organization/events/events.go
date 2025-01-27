package events

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/validator"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

const (
	OrganizationPhoneNumberLinkV1                  = "V1_ORGANIZATION_PHONE_NUMBER_LINK"
	OrganizationRefreshArrV1                       = "V1_ORGANIZATION_REFRESH_ARR"
	OrganizationRefreshRenewalSummaryV1            = "V1_ORGANIZATION_REFRESH_RENEWAL_SUMMARY"
	OrganizationCreateBillingProfileV1             = "V1_ORGANIZATION_CREATE_BILLING_PROFILE"
	OrganizationUpdateBillingProfileV1             = "V1_ORGANIZATION_UPDATE_BILLING_PROFILE"
	OrganizationEmailLinkToBillingProfileV1        = "V1_ORGANIZATION_EMAIL_LINK_TO_BILLING_PROFILE"
	OrganizationEmailUnlinkFromBillingProfileV1    = "V1_ORGANIZATION_EMAIL_UNLINK_FROM_BILLING_PROFILE"
	OrganizationLocationLinkToBillingProfileV1     = "V1_ORGANIZATION_LOCATION_LINK_TO_BILLING_PROFILE"
	OrganizationLocationUnlinkFromBillingProfileV1 = "V1_ORGANIZATION_LOCATION_UNLINK_FROM_BILLING_PROFILE"
)

type OrganizationLinkPhoneNumberEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PhoneNumberId string    `json:"phoneNumberId" validate:"required"`
	Label         string    `json:"label"`
	Primary       bool      `json:"primary"`
}

func NewOrganizationLinkPhoneNumberEvent(aggregate eventstore.Aggregate, phoneNumberId, label string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationLinkPhoneNumberEvent{
		Tenant:        aggregate.GetTenant(),
		UpdatedAt:     updatedAt,
		PhoneNumberId: phoneNumberId,
		Label:         label,
		Primary:       primary,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationLinkPhoneNumberEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationPhoneNumberLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationLinkPhoneNumberEvent")
	}
	return event, nil
}

type OrganizationRefreshArrEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewOrganizationRefreshArrEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := OrganizationRefreshArrEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationRefreshArrEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationRefreshArrV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationRefreshArrEvent")
	}
	return event, nil
}

type OrganizationRefreshRenewalSummaryEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewOrganizationRefreshRenewalSummaryEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := OrganizationRefreshRenewalSummaryEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationRefreshRenewalSummaryEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationRefreshRenewalSummaryV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationRefreshRenewalSummaryEvent")
	}
	return event, nil
}
