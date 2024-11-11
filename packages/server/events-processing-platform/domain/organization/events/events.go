package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"time"

	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

const (
	OrganizationCreateV1          = "V1_ORGANIZATION_CREATE"
	OrganizationUpdateV1          = "V1_ORGANIZATION_UPDATE"
	OrganizationPhoneNumberLinkV1 = "V1_ORGANIZATION_PHONE_NUMBER_LINK"
	// Deprecated
	OrganizationEmailLinkV1 = "V1_ORGANIZATION_EMAIL_LINK"
	// Deprecated
	OrganizationEmailUnlinkV1 = "V1_ORGANIZATION_EMAIL_UNLINK"
	//Deprecated
	OrganizationLocationLinkV1 = "V1_ORGANIZATION_LOCATION_LINK"
	//Deprecated
	OrganizationLinkDomainV1   = "V1_ORGANIZATION_LINK_DOMAIN"
	OrganizationUnlinkDomainV1 = "V1_ORGANIZATION_UNLINK_DOMAIN"
	//Deprecated
	OrganizationAddSocialV1    = "V1_ORGANIZATION_ADD_SOCIAL"
	OrganizationRemoveSocialV1 = "V1_ORGANIZATION_REMOVE_SOCIAL"
	//Deprecated
	OrganizationUpdateRenewalLikelihoodV1 = "V1_ORGANIZATION_UPDATE_RENEWAL_LIKELIHOOD"
	//Deprecated
	OrganizationUpdateRenewalForecastV1 = "V1_ORGANIZATION_UPDATE_RENEWAL_FORECAST"
	//Deprecated
	OrganizationUpdateBillingDetailsV1 = "V1_ORGANIZATION_UPDATE_BILLING_DETAILS"
	//Deprecated
	OrganizationRequestRenewalForecastV1 = "V1_ORGANIZATION_RECALCULATE_RENEWAL_FORECAST_REQUEST"
	//Deprecated
	OrganizationRequestNextCycleDateV1 = "V1_ORGANIZATION_RECALCULATE_NEXT_CYCLE_DATE_REQUEST"
	//Deprecated
	OrganizationRequestScrapeByWebsiteV1 = "V1_ORGANIZATION_SCRAPE_BY_WEBSITE_REQUEST"
	//Deprecated
	OrganizationHideV1 = "V1_ORGANIZATION_HIDE"
	OrganizationShowV1 = "V1_ORGANIZATION_SHOW"
	//Deprecated
	OrganizationRefreshLastTouchpointV1            = "V1_ORGANIZATION_REFRESH_LAST_TOUCHPOINT"
	OrganizationUpsertCustomFieldV1                = "V1_ORGANIZATION_UPSERT_CUSTOM_FIELD"
	OrganizationAddParentV1                        = "V1_ORGANIZATION_ADD_PARENT"
	OrganizationRemoveParentV1                     = "V1_ORGANIZATION_REMOVE_PARENT"
	OrganizationRefreshArrV1                       = "V1_ORGANIZATION_REFRESH_ARR"
	OrganizationRefreshRenewalSummaryV1            = "V1_ORGANIZATION_REFRESH_RENEWAL_SUMMARY"
	OrganizationUpdateOnboardingStatusV1           = "V1_ORGANIZATION_UPDATE_ONBOARDING_STATUS"
	OrganizationUpdateOwnerNotificationV1          = "V1_ORGANIZATION_UPDATE_OWNER_NOTIFICATION"
	OrganizationUpdateOwnerV1                      = "V1_ORGANIZATION_UPDATE_OWNER" //DEPRECATED
	OrganizationCreateBillingProfileV1             = "V1_ORGANIZATION_CREATE_BILLING_PROFILE"
	OrganizationUpdateBillingProfileV1             = "V1_ORGANIZATION_UPDATE_BILLING_PROFILE"
	OrganizationEmailLinkToBillingProfileV1        = "V1_ORGANIZATION_EMAIL_LINK_TO_BILLING_PROFILE"
	OrganizationEmailUnlinkFromBillingProfileV1    = "V1_ORGANIZATION_EMAIL_UNLINK_FROM_BILLING_PROFILE"
	OrganizationLocationLinkToBillingProfileV1     = "V1_ORGANIZATION_LOCATION_LINK_TO_BILLING_PROFILE"
	OrganizationLocationUnlinkFromBillingProfileV1 = "V1_ORGANIZATION_LOCATION_UNLINK_FROM_BILLING_PROFILE"
	OrganizationRequestEnrichV1                    = "V1_ORGANIZATION_ENRICH"
	OrganizationRefreshDerivedDataV1               = "V1_ORGANIZATION_REFRESH_DERIVED_DATA"
	//Deprecated
	OrganizationAddTagV1 = "V1_ORGANIZATION_ADD_TAG"
	//Deprecated
	OrganizationRemoveTagV1      = "V1_ORGANIZATION_REMOVE_TAG"
	OrganizationAddLocationV1    = "V1_ORGANIZATION_ADD_LOCATION"
	OrganizationAdjustIndustryV1 = "V1_ORGANIZATION_ADJUST_INDUSTRY"
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

type OrganizationLinkLocationEvent struct {
	Tenant     string    `json:"tenant" validate:"required"`
	UpdatedAt  time.Time `json:"updatedAt"`
	LocationId string    `json:"locationId" validate:"required"`
}

func NewOrganizationLinkLocationEvent(aggregate eventstore.Aggregate, locationId string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationLinkLocationEvent{
		Tenant:     aggregate.GetTenant(),
		UpdatedAt:  updatedAt,
		LocationId: locationId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationLinkLocationEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationLocationLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationLinkLocationEvent")
	}
	return event, nil
}

type OrganizationRequestScrapeByWebsite struct {
	Tenant      string    `json:"tenant" validate:"required"`
	Website     string    `json:"website" validate:"required"`
	RequestedAt time.Time `json:"requestedAt"`
}

func NewOrganizationRequestScrapeByWebsite(aggregate eventstore.Aggregate, website string) (eventstore.Event, error) {
	eventData := OrganizationRequestScrapeByWebsite{
		Tenant:      aggregate.GetTenant(),
		Website:     website,
		RequestedAt: utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationRequestScrapeByWebsite")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationRequestScrapeByWebsiteV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationRequestScrapeByWebsite")
	}
	return event, nil
}

type ShowOrganizationEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewShowOrganizationEventEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := ShowOrganizationEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ShowOrganizationEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationShowV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ShowOrganizationEvent")
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

type OrganizationAddParentEvent struct {
	Tenant               string `json:"tenant" validate:"required"`
	ParentOrganizationId string `json:"parentOrganizationId" validate:"required"`
	Type                 string `json:"type"`
}

func NewOrganizationAddParentEvent(aggregate eventstore.Aggregate, parentOrganizationId, relType string) (eventstore.Event, error) {
	eventData := OrganizationAddParentEvent{
		Tenant:               aggregate.GetTenant(),
		ParentOrganizationId: parentOrganizationId,
		Type:                 relType,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationAddParentEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationAddParentV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationAddParentEvent")
	}
	return event, nil
}

type OrganizationRemoveParentEvent struct {
	Tenant               string `json:"tenant" validate:"required"`
	ParentOrganizationId string `json:"parentOrganizationId" validate:"required"`
}

func NewOrganizationRemoveParentEvent(aggregate eventstore.Aggregate, parentOrganizationId string) (eventstore.Event, error) {
	eventData := OrganizationRemoveParentEvent{
		Tenant:               aggregate.GetTenant(),
		ParentOrganizationId: parentOrganizationId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationRemoveParentEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationRemoveParentV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationRemoveParentEvent")
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
