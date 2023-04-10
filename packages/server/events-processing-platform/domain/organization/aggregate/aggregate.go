package aggregate

import (
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
	*eventstore.AggregateBase
	Organization *models.Organization
}

func NewOrganizationAggregateWithTenantAndID(tenant, id string) *OrganizationAggregate {
	if id == "" {
		return nil
	}
	aggregate := NewOrganizationAggregate()
	aggregate.SetID(tenant + "-" + id)
	return aggregate
}

func NewOrganizationAggregate() *OrganizationAggregate {
	organizationAggregate := &OrganizationAggregate{Organization: models.NewOrganization()}
	base := eventstore.NewAggregateBase(organizationAggregate.When)
	base.SetType(OrganizationAggregateType)
	organizationAggregate.AggregateBase = base
	return organizationAggregate
}

func (organizationAggregate *OrganizationAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {

	case events.OrganizationCreatedV1:
		return organizationAggregate.onOrganizationCreated(event)
	case events.OrganizationUpdatedV1:
		return organizationAggregate.onOrganizationUpdated(event)
	case events.OrganizationPhoneNumberLinkedV1:
		return organizationAggregate.onPhoneNumberLinked(event)
	case events.OrganizationEmailLinkedV1:
		return organizationAggregate.onEmailLinked(event)

	default:
		return eventstore.ErrInvalidEventType
	}
}

func (a *OrganizationAggregate) onOrganizationCreated(event eventstore.Event) error {
	var eventData events.OrganizationCreatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.Name = eventData.Name
	a.Organization.Description = eventData.Description
	a.Organization.Website = eventData.Website
	a.Organization.Industry = eventData.Industry
	a.Organization.IsPublic = eventData.IsPublic
	a.Organization.Source = common_models.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.Organization.CreatedAt = eventData.CreatedAt
	a.Organization.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *OrganizationAggregate) onOrganizationUpdated(event eventstore.Event) error {
	var eventData events.OrganizationUpdatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.Source.SourceOfTruth = eventData.SourceOfTruth
	a.Organization.UpdatedAt = eventData.UpdatedAt
	a.Organization.Name = eventData.Name
	a.Organization.Description = eventData.Description
	a.Organization.Website = eventData.Website
	a.Organization.Industry = eventData.Industry
	a.Organization.IsPublic = eventData.IsPublic
	return nil
}

func (a *OrganizationAggregate) onPhoneNumberLinked(event eventstore.Event) error {
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

func (a *OrganizationAggregate) onEmailLinked(event eventstore.Event) error {
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
