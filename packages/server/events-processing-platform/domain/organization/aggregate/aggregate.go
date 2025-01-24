package aggregate

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	organizationEvents "github.com/customeros/customeros/packages/server/events-processing-platform/domain/organization/events"
	"github.com/customeros/customeros/packages/server/events-processing-platform/domain/organization/model"
	"github.com/customeros/customeros/packages/server/events/event/common"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	OrganizationAggregateType eventstore.AggregateType = "organization"
)

type OrganizationAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Organization *model.Organization
}

func NewOrganizationAggregateWithTenantAndID(tenant, id string) *OrganizationAggregate {
	organizationAggregate := OrganizationAggregate{}
	organizationAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(OrganizationAggregateType, tenant, id)
	organizationAggregate.SetWhen(organizationAggregate.When)
	organizationAggregate.Organization = &model.Organization{}
	organizationAggregate.Tenant = tenant

	return &organizationAggregate
}

func (a *OrganizationAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.HandleGRPCRequest")
	defer span.Finish()

	return nil, nil
}

func (a *OrganizationAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {
	case organizationEvents.OrganizationPhoneNumberLinkV1:
		return a.onPhoneNumberLink(event)
	case organizationEvents.OrganizationUpsertCustomFieldV1:
		return a.onUpsertCustomField(event)
	case organizationEvents.OrganizationCreateBillingProfileV1:
		return a.onCreateBillingProfile(event)
	case organizationEvents.OrganizationUpdateBillingProfileV1:
		return a.onUpdateBillingProfile(event)
	case organizationEvents.OrganizationEmailLinkToBillingProfileV1:
		return a.onEmailLinkToBillingProfile(event)
	case organizationEvents.OrganizationEmailUnlinkFromBillingProfileV1:
		return a.onEmailUnlinkFromBillingProfile(event)
	case organizationEvents.OrganizationLocationLinkToBillingProfileV1:
		return a.onLocationLinkToBillingProfile(event)
	case organizationEvents.OrganizationLocationUnlinkFromBillingProfileV1:
		return a.onLocationUnlinkFromBillingProfile(event)
	default:
		return nil
	}
}

func (a *OrganizationAggregate) onPhoneNumberLink(event eventstore.Event) error {
	var eventData organizationEvents.OrganizationLinkPhoneNumberEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.PhoneNumbers == nil {
		a.Organization.PhoneNumbers = make(map[string]model.OrganizationPhoneNumber)
	}
	a.Organization.PhoneNumbers[eventData.PhoneNumberId] = model.OrganizationPhoneNumber{
		Label:   eventData.Label,
		Primary: eventData.Primary,
	}
	a.Organization.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *OrganizationAggregate) onUpsertCustomField(event eventstore.Event) error {
	var eventData organizationEvents.OrganizationUpsertCustomField
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.Organization.CustomFields == nil {
		a.Organization.CustomFields = make(map[string]model.CustomField)
	}

	if val, ok := a.Organization.CustomFields[eventData.CustomFieldId]; ok {
		val.Source.SourceOfTruth = eventData.SourceOfTruth
		val.UpdatedAt = eventData.UpdatedAt
		val.CustomFieldValue = eventData.CustomFieldValue
		val.Name = eventData.CustomFieldName
	} else {
		a.Organization.CustomFields[eventData.CustomFieldId] = model.CustomField{
			Source: common.Source{
				Source:        eventData.Source,
				SourceOfTruth: eventData.SourceOfTruth,
				AppSource:     eventData.AppSource,
			},
			CreatedAt:           eventData.CreatedAt,
			UpdatedAt:           eventData.UpdatedAt,
			Id:                  eventData.CustomFieldId,
			TemplateId:          eventData.TemplateId,
			Name:                eventData.CustomFieldName,
			CustomFieldDataType: model.CustomFieldDataType(eventData.CustomFieldDataType),
			CustomFieldValue:    eventData.CustomFieldValue,
		}
	}
	return nil
}

func (a *OrganizationAggregate) onCreateBillingProfile(event eventstore.Event) error {
	var eventData organizationEvents.BillingProfileCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.Organization.BillingProfiles == nil {
		a.Organization.BillingProfiles = make(map[string]model.BillingProfile)
	}

	a.Organization.BillingProfiles[eventData.BillingProfileId] = model.BillingProfile{
		Id:           eventData.BillingProfileId,
		LegalName:    eventData.LegalName,
		TaxId:        eventData.TaxId,
		CreatedAt:    eventData.CreatedAt,
		UpdatedAt:    eventData.UpdatedAt,
		SourceFields: eventData.SourceFields,
	}

	return nil
}

func (a *OrganizationAggregate) onUpdateBillingProfile(event eventstore.Event) error {
	var eventData organizationEvents.BillingProfileUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.Organization.BillingProfiles == nil {
		a.Organization.BillingProfiles = make(map[string]model.BillingProfile)
	}
	billingProfile, ok := a.Organization.BillingProfiles[eventData.BillingProfileId]
	if !ok {
		a.Organization.BillingProfiles[eventData.BillingProfileId] = model.BillingProfile{}
		billingProfile = a.Organization.BillingProfiles[eventData.BillingProfileId]
	}

	if eventData.UpdateLegalName() {
		billingProfile.LegalName = eventData.LegalName
	}
	if eventData.UpdateTaxId() {
		billingProfile.TaxId = eventData.TaxId
	}
	billingProfile.UpdatedAt = eventData.UpdatedAt
	a.Organization.BillingProfiles[eventData.BillingProfileId] = billingProfile

	return nil
}

func (a *OrganizationAggregate) onEmailLinkToBillingProfile(event eventstore.Event) error {
	var eventData organizationEvents.LinkEmailToBillingProfileEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.BillingProfiles == nil {
		a.Organization.BillingProfiles = make(map[string]model.BillingProfile)
	}
	billingProfile, ok := a.Organization.BillingProfiles[eventData.BillingProfileId]
	if !ok {
		a.Organization.BillingProfiles[eventData.BillingProfileId] = model.BillingProfile{}
		billingProfile = a.Organization.BillingProfiles[eventData.BillingProfileId]
	}

	if eventData.Primary {
		billingProfile.PrimaryEmailId = eventData.EmailId
		billingProfile.EmailIds = utils.RemoveFromList(billingProfile.EmailIds, eventData.EmailId)
	} else {
		billingProfile.EmailIds = utils.AddToListIfNotExists(billingProfile.EmailIds, eventData.EmailId)
		if billingProfile.PrimaryEmailId == eventData.EmailId {
			billingProfile.PrimaryEmailId = ""
		}
	}
	billingProfile.UpdatedAt = eventData.UpdatedAt
	a.Organization.BillingProfiles[eventData.BillingProfileId] = billingProfile

	return nil
}

func (a *OrganizationAggregate) onEmailUnlinkFromBillingProfile(event eventstore.Event) error {
	var eventData organizationEvents.UnlinkEmailFromBillingProfileEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.BillingProfiles == nil {
		a.Organization.BillingProfiles = make(map[string]model.BillingProfile)
	}
	billingProfile, ok := a.Organization.BillingProfiles[eventData.BillingProfileId]
	if !ok {
		a.Organization.BillingProfiles[eventData.BillingProfileId] = model.BillingProfile{}
		billingProfile = a.Organization.BillingProfiles[eventData.BillingProfileId]
	}
	if billingProfile.PrimaryEmailId == eventData.EmailId {
		billingProfile.PrimaryEmailId = ""
	}
	billingProfile.EmailIds = utils.RemoveFromList(billingProfile.EmailIds, eventData.EmailId)
	billingProfile.UpdatedAt = eventData.UpdatedAt
	a.Organization.BillingProfiles[eventData.BillingProfileId] = billingProfile

	return nil
}

func (a *OrganizationAggregate) onLocationLinkToBillingProfile(event eventstore.Event) error {
	var eventData organizationEvents.LinkLocationToBillingProfileEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.BillingProfiles == nil {
		a.Organization.BillingProfiles = make(map[string]model.BillingProfile)
	}
	billingProfile, ok := a.Organization.BillingProfiles[eventData.BillingProfileId]
	if !ok {
		a.Organization.BillingProfiles[eventData.BillingProfileId] = model.BillingProfile{}
		billingProfile = a.Organization.BillingProfiles[eventData.BillingProfileId]
	}

	billingProfile.LocationIds = utils.AddToListIfNotExists(billingProfile.LocationIds, eventData.LocationId)
	billingProfile.UpdatedAt = eventData.UpdatedAt
	a.Organization.BillingProfiles[eventData.BillingProfileId] = billingProfile

	return nil
}

func (a *OrganizationAggregate) onLocationUnlinkFromBillingProfile(event eventstore.Event) error {
	var eventData organizationEvents.UnlinkLocationFromBillingProfileEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.BillingProfiles == nil {
		a.Organization.BillingProfiles = make(map[string]model.BillingProfile)
	}
	billingProfile, ok := a.Organization.BillingProfiles[eventData.BillingProfileId]
	if !ok {
		a.Organization.BillingProfiles[eventData.BillingProfileId] = model.BillingProfile{}
		billingProfile = a.Organization.BillingProfiles[eventData.BillingProfileId]
	}
	billingProfile.LocationIds = utils.RemoveFromList(billingProfile.LocationIds, eventData.LocationId)
	billingProfile.UpdatedAt = eventData.UpdatedAt
	a.Organization.BillingProfiles[eventData.BillingProfileId] = billingProfile

	return nil
}
