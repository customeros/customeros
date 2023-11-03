package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"strings"
)

const (
	OrganizationAggregateType eventstore.AggregateType = "organization"
)

type OrganizationAggregate struct {
	*aggregate.CommonTenantIdAggregate
	Organization *models.Organization
}

func NewOrganizationAggregateWithTenantAndID(tenant, id string) *OrganizationAggregate {
	organizationAggregate := OrganizationAggregate{}
	organizationAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(OrganizationAggregateType, tenant, id)
	organizationAggregate.SetWhen(organizationAggregate.When)
	organizationAggregate.Organization = &models.Organization{}
	organizationAggregate.Tenant = tenant

	return &organizationAggregate
}

func (a *OrganizationAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {
	case events.OrganizationCreateV1:
		return a.onOrganizationCreate(event)
	case events.OrganizationUpdateV1:
		return a.onOrganizationUpdate(event)
	case events.OrganizationPhoneNumberLinkV1:
		return a.onPhoneNumberLink(event)
	case events.OrganizationEmailLinkV1:
		return a.onEmailLink(event)
	case events.OrganizationLocationLinkV1:
		return a.onLocationLink(event)
	case events.OrganizationLinkDomainV1:
		return a.onDomainLink(event)
	case events.OrganizationAddSocialV1:
		return a.onAddSocial(event)
	case events.OrganizationUpdateRenewalLikelihoodV1:
		return a.onUpdateRenewalLikelihood(event)
	case events.OrganizationUpdateRenewalForecastV1:
		return a.onUpdateRenewalForecast(event)
	case events.OrganizationUpdateBillingDetailsV1:
		return a.onUpdateBillingDetails(event)
	case events.OrganizationHideV1:
		return a.onHide(event)
	case events.OrganizationShowV1:
		return a.onShow(event)
	case events.OrganizationUpsertCustomFieldV1:
		return a.onUpsertCustomField(event)
	case events.OrganizationAddParentV1:
		return a.onAddParent(event)
	case events.OrganizationRemoveParentV1:
		return a.onRemoveParent(event)
	case events.OrganizationRequestRenewalForecastV1,
		events.OrganizationRequestNextCycleDateV1,
		events.OrganizationRefreshLastTouchpointV1,
		events.OrganizationRequestScrapeByWebsiteV1:
		return nil
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *OrganizationAggregate) onOrganizationCreate(event eventstore.Event) error {
	var eventData events.OrganizationCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.Name = eventData.Name
	a.Organization.Hide = eventData.Hide
	a.Organization.Description = eventData.Description
	a.Organization.Website = eventData.Website
	a.Organization.Industry = eventData.Industry
	a.Organization.SubIndustry = eventData.SubIndustry
	a.Organization.IndustryGroup = eventData.IndustryGroup
	a.Organization.TargetAudience = eventData.TargetAudience
	a.Organization.ValueProposition = eventData.ValueProposition
	a.Organization.LastFundingRound = eventData.LastFundingRound
	a.Organization.LastFundingAmount = eventData.LastFundingAmount
	a.Organization.ReferenceId = eventData.ReferenceId
	a.Organization.Note = eventData.Note
	a.Organization.IsPublic = eventData.IsPublic
	a.Organization.IsCustomer = eventData.IsCustomer
	a.Organization.Employees = eventData.Employees
	a.Organization.Market = eventData.Market
	a.Organization.Source = cmnmod.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.Organization.CreatedAt = eventData.CreatedAt
	a.Organization.UpdatedAt = eventData.UpdatedAt
	if eventData.ExternalSystem.Available() {
		a.Organization.ExternalSystems = []cmnmod.ExternalSystem{eventData.ExternalSystem}
	}
	return nil
}

func (a *OrganizationAggregate) onOrganizationUpdate(event eventstore.Event) error {
	var eventData events.OrganizationUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if eventData.Source == constants.SourceOpenline {
		a.Organization.Source.SourceOfTruth = eventData.Source
	}
	a.Organization.UpdatedAt = eventData.UpdatedAt

	if eventData.Source != a.Organization.Source.SourceOfTruth && a.Organization.Source.SourceOfTruth == constants.SourceOpenline {
		if a.Organization.Name == "" {
			a.Organization.Name = eventData.Name
		}
		if a.Organization.Description == "" {
			a.Organization.Description = eventData.Description
		}
		if a.Organization.Website == "" {
			a.Organization.Website = eventData.Website
		}
		if a.Organization.Industry == "" {
			a.Organization.Industry = eventData.Industry
		}
		if a.Organization.SubIndustry == "" {
			a.Organization.SubIndustry = eventData.SubIndustry
		}
		if a.Organization.IndustryGroup == "" {
			a.Organization.IndustryGroup = eventData.IndustryGroup
		}
		if a.Organization.TargetAudience == "" {
			a.Organization.TargetAudience = eventData.TargetAudience
		}
		if a.Organization.ValueProposition == "" {
			a.Organization.ValueProposition = eventData.ValueProposition
		}
		if a.Organization.LastFundingRound == "" {
			a.Organization.LastFundingRound = eventData.LastFundingRound
		}
		if a.Organization.LastFundingAmount == "" {
			a.Organization.LastFundingAmount = eventData.LastFundingAmount
		}
		if a.Organization.ReferenceId == "" {
			a.Organization.ReferenceId = eventData.ReferenceId
		}
		if a.Organization.Note == "" {
			a.Organization.Note = eventData.Note
		}
		if a.Organization.Employees == 0 {
			a.Organization.Employees = eventData.Employees
		}
		if a.Organization.Market == "" {
			a.Organization.Market = eventData.Market
		}
		if !a.Organization.IsCustomer {
			a.Organization.IsCustomer = eventData.IsCustomer
		}
	} else {
		if !eventData.IgnoreEmptyFields {
			a.Organization.IsPublic = eventData.IsPublic
			a.Organization.IsCustomer = eventData.IsCustomer
			a.Organization.Hide = eventData.Hide
			a.Organization.Name = eventData.Name
			a.Organization.Description = eventData.Description
			a.Organization.Website = eventData.Website
			a.Organization.Industry = eventData.Industry
			a.Organization.SubIndustry = eventData.SubIndustry
			a.Organization.IndustryGroup = eventData.IndustryGroup
			a.Organization.TargetAudience = eventData.TargetAudience
			a.Organization.ValueProposition = eventData.ValueProposition
			a.Organization.LastFundingRound = eventData.LastFundingRound
			a.Organization.LastFundingAmount = eventData.LastFundingAmount
			a.Organization.ReferenceId = eventData.ReferenceId
			a.Organization.Note = eventData.Note
			a.Organization.Employees = eventData.Employees
			a.Organization.Market = eventData.Market
		} else {
			if eventData.Name != "" {
				a.Organization.Name = eventData.Name
			}
			if eventData.Description != "" {
				a.Organization.Description = eventData.Description
			}
			if eventData.Website != "" {
				a.Organization.Website = eventData.Website
			}
			if eventData.Industry != "" {
				a.Organization.Industry = eventData.Industry
			}
			if eventData.SubIndustry != "" {
				a.Organization.SubIndustry = eventData.SubIndustry
			}
			if eventData.IndustryGroup != "" {
				a.Organization.IndustryGroup = eventData.IndustryGroup
			}
			if eventData.TargetAudience != "" {
				a.Organization.TargetAudience = eventData.TargetAudience
			}
			if eventData.ValueProposition != "" {
				a.Organization.ValueProposition = eventData.ValueProposition
			}
			if eventData.LastFundingRound != "" {
				a.Organization.LastFundingRound = eventData.LastFundingRound
			}
			if eventData.LastFundingAmount != "" {
				a.Organization.LastFundingAmount = eventData.LastFundingAmount
			}
			if eventData.ReferenceId != "" {
				a.Organization.ReferenceId = eventData.ReferenceId
			}
			if eventData.Note != "" {
				a.Organization.Note = eventData.Note
			}
			if eventData.Employees != 0 {
				a.Organization.Employees = eventData.Employees
			}
			if eventData.Market != "" {
				a.Organization.Market = eventData.Market
			}
			if eventData.IsCustomer {
				a.Organization.IsCustomer = eventData.IsCustomer
			}
		}
	}
	if eventData.ExternalSystem.Available() {
		found := false
		for _, externalSystem := range a.Organization.ExternalSystems {
			if externalSystem.ExternalSystemId == eventData.ExternalSystem.ExternalSystemId &&
				externalSystem.ExternalId == eventData.ExternalSystem.ExternalId {
				found = true
				externalSystem.ExternalUrl = eventData.ExternalSystem.ExternalUrl
				externalSystem.SyncDate = eventData.ExternalSystem.SyncDate
				externalSystem.ExternalSource = eventData.ExternalSystem.ExternalSource
				if eventData.ExternalSystem.ExternalIdSecond != "" {
					externalSystem.ExternalIdSecond = eventData.ExternalSystem.ExternalIdSecond
				}
			}
		}
		if !found {
			a.Organization.ExternalSystems = append(a.Organization.ExternalSystems, eventData.ExternalSystem)
		}
	}
	return nil
}

func (a *OrganizationAggregate) onPhoneNumberLink(event eventstore.Event) error {
	var eventData events.OrganizationLinkPhoneNumberEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.PhoneNumbers == nil {
		a.Organization.PhoneNumbers = make(map[string]models.OrganizationPhoneNumber)
	}
	a.Organization.PhoneNumbers[eventData.PhoneNumberId] = models.OrganizationPhoneNumber{
		Label:   eventData.Label,
		Primary: eventData.Primary,
	}
	a.Organization.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *OrganizationAggregate) onEmailLink(event eventstore.Event) error {
	var eventData events.OrganizationLinkEmailEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.Emails == nil {
		a.Organization.Emails = make(map[string]models.OrganizationEmail)
	}
	a.Organization.Emails[eventData.EmailId] = models.OrganizationEmail{
		Label:   eventData.Label,
		Primary: eventData.Primary,
	}
	a.Organization.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *OrganizationAggregate) onDomainLink(event eventstore.Event) error {
	var eventData events.OrganizationLinkDomainEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.Domains == nil {
		a.Organization.Domains = []string{}
	}
	if !utils.Contains(a.Organization.Domains, strings.TrimSpace(eventData.Domain)) {
		a.Organization.Domains = append(a.Organization.Domains, strings.TrimSpace(eventData.Domain))
	}
	return nil
}

func (a *OrganizationAggregate) onLocationLink(event eventstore.Event) error {
	var eventData events.OrganizationLinkLocationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.Locations = utils.AddToListIfNotExists(a.Organization.Locations, eventData.LocationId)
	return nil
}

func (a *OrganizationAggregate) onAddSocial(event eventstore.Event) error {
	var eventData events.OrganizationAddSocialEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.Socials == nil {
		a.Organization.Socials = make(map[string]models.Social)
	}
	a.Organization.Socials[eventData.SocialId] = models.Social{
		PlatformName: eventData.PlatformName,
		Url:          eventData.Url,
	}
	return nil
}

func (a *OrganizationAggregate) onUpdateRenewalLikelihood(event eventstore.Event) error {
	var eventData events.OrganizationUpdateRenewalLikelihoodEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.RenewalLikelihood.RenewalLikelihood = eventData.RenewalLikelihood
	a.Organization.RenewalLikelihood.Comment = eventData.Comment
	a.Organization.RenewalLikelihood.UpdatedBy = eventData.UpdatedBy
	a.Organization.RenewalLikelihood.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *OrganizationAggregate) onUpdateRenewalForecast(event eventstore.Event) error {
	var eventData events.OrganizationUpdateRenewalForecastEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.RenewalForecast.Amount = eventData.Amount
	a.Organization.RenewalForecast.Comment = eventData.Comment
	a.Organization.RenewalForecast.UpdatedBy = eventData.UpdatedBy
	a.Organization.RenewalForecast.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *OrganizationAggregate) onUpdateBillingDetails(event eventstore.Event) error {
	var eventData events.OrganizationUpdateBillingDetailsEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.BillingDetails.Amount = eventData.Amount
	a.Organization.BillingDetails.Frequency = eventData.Frequency
	a.Organization.BillingDetails.RenewalCycle = eventData.RenewalCycle
	a.Organization.BillingDetails.RenewalCycleStart = eventData.RenewalCycleStart
	if eventData.UpdatedBy == "" {
		a.Organization.BillingDetails.RenewalCycleNext = eventData.RenewalCycleNext
	}
	return nil
}

func (a *OrganizationAggregate) onHide(event eventstore.Event) error {
	var eventData events.HideOrganizationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.Hide = true
	return nil
}

func (a *OrganizationAggregate) onShow(event eventstore.Event) error {
	var eventData events.ShowOrganizationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.Hide = false
	return nil
}

func (a *OrganizationAggregate) onUpsertCustomField(event eventstore.Event) error {
	var eventData events.OrganizationUpsertCustomField
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.Organization.CustomFields == nil {
		a.Organization.CustomFields = make(map[string]models.CustomField)
	}

	if val, ok := a.Organization.CustomFields[eventData.CustomFieldId]; ok {
		val.Source.SourceOfTruth = eventData.SourceOfTruth
		val.UpdatedAt = eventData.UpdatedAt
		val.CustomFieldValue = eventData.CustomFieldValue
		val.Name = eventData.CustomFieldName
	} else {
		a.Organization.CustomFields[eventData.CustomFieldId] = models.CustomField{
			Source: cmnmod.Source{
				Source:        eventData.Source,
				SourceOfTruth: eventData.SourceOfTruth,
				AppSource:     eventData.AppSource,
			},
			CreatedAt:           eventData.CreatedAt,
			UpdatedAt:           eventData.UpdatedAt,
			Id:                  eventData.CustomFieldId,
			TemplateId:          eventData.TemplateId,
			Name:                eventData.CustomFieldName,
			CustomFieldDataType: models.CustomFieldDataType(eventData.CustomFieldDataType),
			CustomFieldValue:    eventData.CustomFieldValue,
		}
	}
	return nil
}

func (a *OrganizationAggregate) onAddParent(event eventstore.Event) error {
	var eventData events.OrganizationAddParentEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.ParentOrganizations == nil {
		a.Organization.ParentOrganizations = make(map[string]models.ParentOrganization)
	}
	a.Organization.ParentOrganizations[eventData.ParentOrganizationId] = models.ParentOrganization{
		OrganizationId: eventData.ParentOrganizationId,
		Type:           eventData.Type,
	}
	return nil
}

func (a *OrganizationAggregate) onRemoveParent(event eventstore.Event) error {
	var eventData events.OrganizationRemoveParentEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	delete(a.Organization.ParentOrganizations, eventData.ParentOrganizationId)
	return nil
}
