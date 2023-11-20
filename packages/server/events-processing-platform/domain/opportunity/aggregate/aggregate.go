package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
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
	case event.OpportunityUpdateV1:
		return a.onOpportunityUpdate(evt)
	case event.OpportunityCreateRenewalV1:
		return a.onRenewalOpportunityCreate(evt)
	case event.OpportunityUpdateRenewalV1:
		return a.onRenewalOpportunityUpdate(evt)
	case event.OpportunityUpdateNextCycleDateV1:
		return a.onOpportunityUpdateNextCycleDate(evt)
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

func (a *OpportunityAggregate) onRenewalOpportunityCreate(evt eventstore.Event) error {
	var eventData event.OpportunityCreateRenewalEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Opportunity.ID = a.ID
	a.Opportunity.Tenant = a.Tenant
	a.Opportunity.ContractId = eventData.ContractId
	a.Opportunity.InternalType = model.OpportunityInternalTypeStringRenewal
	a.Opportunity.InternalStage = model.OpportunityInternalStageStringDecode(eventData.InternalStage)
	a.Opportunity.CreatedAt = eventData.CreatedAt
	a.Opportunity.UpdatedAt = eventData.UpdatedAt
	a.Opportunity.Source = eventData.Source
	a.Opportunity.RenewalDetails.RenewalLikelihood = eventData.RenewalLikelihood

	return nil
}

func (a *OpportunityAggregate) onOpportunityUpdateNextCycleDate(evt eventstore.Event) error {
	var eventData event.OpportunityUpdateNextCycleDateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Opportunity.RenewalDetails.RenewedAt = eventData.RenewedAt

	return nil
}

func (a *OpportunityAggregate) onOpportunityUpdate(evt eventstore.Event) error {
	var eventData event.OpportunityUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	// Update only if the source of truth is 'openline' or the new source matches the source of truth
	if eventData.Source == constants.SourceOpenline {
		a.Opportunity.Source.SourceOfTruth = eventData.Source
	}

	if eventData.Source != a.Opportunity.Source.SourceOfTruth && a.Opportunity.Source.SourceOfTruth == constants.SourceOpenline {
		// Update fields only if they are empty
		if a.Opportunity.Name == "" && eventData.UpdateName() {
			a.Opportunity.Name = eventData.Name
		}
	} else {
		if eventData.UpdateName() {
			a.Opportunity.Name = eventData.Name
		}
		if eventData.UpdateAmount() {
			a.Opportunity.Amount = eventData.Amount
		}
		if eventData.UpdateMaxAmount() {
			a.Opportunity.MaxAmount = eventData.MaxAmount
		}
	}

	a.Opportunity.UpdatedAt = eventData.UpdatedAt

	if eventData.ExternalSystem.Available() {
		found := false
		for _, externalSystem := range a.Opportunity.ExternalSystems {
			if externalSystem.ExternalSystemId == eventData.ExternalSystem.ExternalSystemId && externalSystem.ExternalId == eventData.ExternalSystem.ExternalId {
				found = true
				externalSystem.ExternalUrl = eventData.ExternalSystem.ExternalUrl
				externalSystem.ExternalSource = eventData.ExternalSystem.ExternalSource
				externalSystem.SyncDate = eventData.ExternalSystem.SyncDate
				if eventData.ExternalSystem.ExternalIdSecond != "" {
					externalSystem.ExternalIdSecond = eventData.ExternalSystem.ExternalIdSecond
				}
			}
		}
		if !found {
			a.Opportunity.ExternalSystems = append(a.Opportunity.ExternalSystems, eventData.ExternalSystem)
		}
	}

	return nil
}

func (a *OpportunityAggregate) onRenewalOpportunityUpdate(evt eventstore.Event) error {
	var eventData event.OpportunityUpdateRenewalEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Opportunity.UpdatedAt = eventData.UpdatedAt
	a.Opportunity.RenewalDetails.RenewalLikelihood = eventData.RenewalLikelihood
	if eventData.UpdatedByUserId != "" &&
		(eventData.Amount != a.Opportunity.Amount || eventData.RenewalLikelihood != a.Opportunity.RenewalDetails.RenewalLikelihood) {
		a.Opportunity.RenewalDetails.RenewalUpdatedByUserAt = &eventData.UpdatedAt
		a.Opportunity.RenewalDetails.RenewalUpdatedByUserId = eventData.UpdatedByUserId
	}
	a.Opportunity.Comments = eventData.Comments
	a.Opportunity.Amount = eventData.Amount
	if eventData.Source == constants.SourceOpenline {
		a.Opportunity.Source.SourceOfTruth = eventData.Source
	}

	return nil
}
