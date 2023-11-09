package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	OpportunityAggregateType eventstore.AggregateType = "opportunity"
)

type OpportunityAggregate struct {
	*aggregate.CommonTenantIdAggregate
	Opportunity *model.Opportunity
}

func NewOpportunityAggregateWithTenantAndID(tenant, id string) *OpportunityAggregate {
	oppAggregate := OpportunityAggregate{}
	oppAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(OpportunityAggregateType, tenant, id)
	oppAggregate.SetWhen(oppAggregate.When)
	oppAggregate.Opportunity = &model.Opportunity{}
	oppAggregate.Tenant = tenant

	return &oppAggregate
}

func (a *OpportunityAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.OpportunityCreateV1:
		return a.onOpportunityCreate(evt)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *OpportunityAggregate) onOpportunityCreate(evt eventstore.Event) error {
	var eventData event.OpportunityCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Opportunity.ID = a.ID
	a.Opportunity.Tenant = a.Tenant
	a.Opportunity.OrganizationId = eventData.OrganizationId
	a.Opportunity.Name = eventData.Name
	a.Opportunity.Amount = eventData.Amount
	a.Opportunity.InternalType = eventData.InternalType
	a.Opportunity.ExternalType = eventData.ExternalType
	a.Opportunity.InternalStage = eventData.InternalStage
	a.Opportunity.ExternalStage = eventData.ExternalStage
	a.Opportunity.EstimatedClosedAt = eventData.EstimatedClosedAt
	a.Opportunity.OwnerUserId = eventData.OwnerUserId
	a.Opportunity.CreatedByUserId = eventData.CreatedByUserId
	a.Opportunity.GeneralNotes = eventData.GeneralNotes
	a.Opportunity.NextSteps = eventData.NextSteps
	a.Opportunity.CreatedAt = eventData.CreatedAt
	a.Opportunity.UpdatedAt = eventData.UpdatedAt
	a.Opportunity.Source = eventData.Source
	if eventData.ExternalSystem.Available() {
		a.Opportunity.ExternalSystems = []commonmodel.ExternalSystem{eventData.ExternalSystem}
	}

	return nil
}
