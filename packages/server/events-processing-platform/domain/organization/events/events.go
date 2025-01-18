package events

import (
	"github.com/customeros/customeros/packages/server/events/event/common"
	"time"

	neo4jmodel "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/validator"
	"github.com/customeros/customeros/packages/server/events-processing-platform/domain/organization/model"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

const (
	OrganizationPhoneNumberLinkV1                  = "V1_ORGANIZATION_PHONE_NUMBER_LINK"
	OrganizationUpsertCustomFieldV1                = "V1_ORGANIZATION_UPSERT_CUSTOM_FIELD"
	OrganizationRefreshArrV1                       = "V1_ORGANIZATION_REFRESH_ARR"
	OrganizationRefreshRenewalSummaryV1            = "V1_ORGANIZATION_REFRESH_RENEWAL_SUMMARY"
	OrganizationUpdateOwnerNotificationV1          = "V1_ORGANIZATION_UPDATE_OWNER_NOTIFICATION"
	OrganizationUpdateOwnerV1                      = "V1_ORGANIZATION_UPDATE_OWNER" //DEPRECATED
	OrganizationCreateBillingProfileV1             = "V1_ORGANIZATION_CREATE_BILLING_PROFILE"
	OrganizationUpdateBillingProfileV1             = "V1_ORGANIZATION_UPDATE_BILLING_PROFILE"
	OrganizationEmailLinkToBillingProfileV1        = "V1_ORGANIZATION_EMAIL_LINK_TO_BILLING_PROFILE"
	OrganizationEmailUnlinkFromBillingProfileV1    = "V1_ORGANIZATION_EMAIL_UNLINK_FROM_BILLING_PROFILE"
	OrganizationLocationLinkToBillingProfileV1     = "V1_ORGANIZATION_LOCATION_LINK_TO_BILLING_PROFILE"
	OrganizationLocationUnlinkFromBillingProfileV1 = "V1_ORGANIZATION_LOCATION_UNLINK_FROM_BILLING_PROFILE"
	OrganizationRefreshDerivedDataV1               = "V1_ORGANIZATION_REFRESH_DERIVED_DATA"
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

type OrganizationUpsertCustomField struct {
	Tenant              string                      `json:"tenant" validate:"required"`
	Source              string                      `json:"source,omitempty"`
	SourceOfTruth       string                      `json:"sourceOfTruth,omitempty"`
	AppSource           string                      `json:"appSource,omitempty"`
	CreatedAt           time.Time                   `json:"createdAt"`
	UpdatedAt           time.Time                   `json:"updatedAt"`
	ExistsInEventStore  bool                        `json:"existsInEventStore"`
	TemplateId          *string                     `json:"templateId,omitempty"`
	CustomFieldId       string                      `json:"customFieldId"`
	CustomFieldName     string                      `json:"customFieldName"`
	CustomFieldDataType string                      `json:"customFieldDataType"`
	CustomFieldValue    neo4jmodel.CustomFieldValue `json:"customFieldValue"`
}

func NewOrganizationUpsertCustomField(aggregate eventstore.Aggregate, sourceFields common.Source, createdAt, updatedAt time.Time, customField model.CustomField, foundInEventStore bool) (eventstore.Event, error) {
	eventData := OrganizationUpsertCustomField{
		Tenant:              aggregate.GetTenant(),
		Source:              sourceFields.Source,
		SourceOfTruth:       sourceFields.SourceOfTruth,
		AppSource:           sourceFields.AppSource,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
		ExistsInEventStore:  foundInEventStore,
		CustomFieldId:       customField.Id,
		TemplateId:          customField.TemplateId,
		CustomFieldName:     customField.Name,
		CustomFieldDataType: string(customField.CustomFieldDataType),
		CustomFieldValue:    customField.CustomFieldValue,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationUpsertCustomField")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationUpsertCustomFieldV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationUpsertCustomField")
	}
	return event, nil
}

type OrganizationOwnerUpdateEvent struct {
	Tenant         string    `json:"tenant" validate:"required"`
	UpdatedAt      time.Time `json:"updatedAt"`
	OwnerUserId    string    `json:"ownerUserId" validate:"required"` // who became owner
	OrganizationId string    `json:"organizationId" validate:"required"`
	ActorUserId    string    `json:"actorUserId"` // who set the owner
}

func NewOrganizationOwnerUpdateEvent(aggregate eventstore.Aggregate, ownerUserId, actorUserId, organizationId string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationOwnerUpdateEvent{
		Tenant:         aggregate.GetTenant(),
		UpdatedAt:      updatedAt,
		OwnerUserId:    ownerUserId,
		OrganizationId: organizationId,
		ActorUserId:    actorUserId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationOwnerUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationUpdateOwnerV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationOwnerUpdateEvent")
	}
	return event, nil
}

func NewOrganizationOwnerUpdateNotificationEvent(aggregate eventstore.Aggregate, ownerUserId, actorUserId, organizationId string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationOwnerUpdateEvent{
		Tenant:         aggregate.GetTenant(),
		UpdatedAt:      updatedAt,
		OwnerUserId:    ownerUserId,
		OrganizationId: organizationId,
		ActorUserId:    actorUserId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationOwnerUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationUpdateOwnerNotificationV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationOwnerUpdateEvent")
	}
	return event, nil
}
