package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
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
	a.Organization.Note = eventData.Note
	a.Organization.IsPublic = eventData.IsPublic
	a.Organization.Employees = eventData.Employees
	a.Organization.Market = eventData.Market
	a.Organization.Source = common_models.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.Organization.CreatedAt = eventData.CreatedAt
	a.Organization.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *OrganizationAggregate) onOrganizationUpdate(event eventstore.Event) error {
	var eventData events.OrganizationUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.Source.SourceOfTruth = eventData.SourceOfTruth
	a.Organization.UpdatedAt = eventData.UpdatedAt
	if !eventData.IgnoreEmptyFields {
		a.Organization.IsPublic = eventData.IsPublic
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
		if eventData.Note != "" {
			a.Organization.Note = eventData.Note
		}
		if eventData.Employees != 0 {
			a.Organization.Employees = eventData.Employees
		}
		if eventData.Market != "" {
			a.Organization.Market = eventData.Market
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
	if !utils.Contains(a.Organization.Domains, eventData.Domain) {
		a.Organization.Domains = append(a.Organization.Domains, eventData.Domain)
	}
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
