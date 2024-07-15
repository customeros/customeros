package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events/utils"
	"github.com/pkg/errors"
	"strings"
)

const (
	MasterPlanAggregateType eventstore.AggregateType = "master_plan"
)

type MasterPlanAggregate struct {
	*eventstore.CommonTenantIdAggregate
	MasterPlan *model.MasterPlan
}

func NewMasterPlanAggregateWithTenantAndID(tenant, id string) *MasterPlanAggregate {
	masterPlanAggregate := MasterPlanAggregate{}
	masterPlanAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(MasterPlanAggregateType, tenant, id)
	masterPlanAggregate.SetWhen(masterPlanAggregate.When)
	masterPlanAggregate.MasterPlan = &model.MasterPlan{}
	masterPlanAggregate.Tenant = tenant

	return &masterPlanAggregate
}

func (a *MasterPlanAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.MasterPlanCreateV1:
		return a.onMasterPlanCreate(evt)
	case event.MasterPlanUpdateV1:
		return a.onMasterPlanUpdate(evt)
	case event.MasterPlanMilestoneCreateV1:
		return a.onMasterPlanMilestoneCreate(evt)
	case event.MasterPlanMilestoneUpdateV1:
		return a.onMasterPlanMilestoneUpdate(evt)
	case event.MasterPlanMilestoneReorderV1:
		return a.onMasterPlanMilestoneReorder(evt)
	default:
		if strings.HasPrefix(evt.GetEventType(), utils.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *MasterPlanAggregate) onMasterPlanCreate(evt eventstore.Event) error {
	var eventData event.MasterPlanCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.MasterPlan.ID = a.ID
	a.MasterPlan.Name = eventData.Name
	a.MasterPlan.CreatedAt = eventData.CreatedAt
	a.MasterPlan.SourceFields = eventData.SourceFields

	return nil
}

func (a *MasterPlanAggregate) onMasterPlanMilestoneCreate(evt eventstore.Event) error {
	var eventData event.MasterPlanMilestoneCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	milestone := model.MasterPlanMilestone{
		ID:            eventData.MilestoneId,
		Name:          eventData.Name,
		Order:         eventData.Order,
		CreatedAt:     eventData.CreatedAt,
		SourceFields:  eventData.SourceFields,
		DurationHours: eventData.DurationHours,
		Items:         eventData.Items,
		Optional:      eventData.Optional,
	}

	if a.MasterPlan.Milestones == nil {
		a.MasterPlan.Milestones = make(map[string]model.MasterPlanMilestone)
	}
	a.MasterPlan.Milestones[milestone.ID] = milestone

	return nil
}

func (a *MasterPlanAggregate) onMasterPlanMilestoneUpdate(evt eventstore.Event) error {
	var eventData event.MasterPlanMilestoneUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.MasterPlan.Milestones == nil {
		a.MasterPlan.Milestones = make(map[string]model.MasterPlanMilestone)
	}
	if _, ok := a.MasterPlan.Milestones[eventData.MilestoneId]; !ok {
		a.MasterPlan.Milestones[eventData.MilestoneId] = model.MasterPlanMilestone{
			ID: eventData.MilestoneId,
		}
	}
	milestone := a.MasterPlan.Milestones[eventData.MilestoneId]
	if eventData.UpdateName() {
		milestone.Name = eventData.Name
	}
	if eventData.UpdateOrder() {
		milestone.Order = eventData.Order
	}
	if eventData.UpdateDurationHours() {
		milestone.DurationHours = eventData.DurationHours
	}
	if eventData.UpdateItems() {
		milestone.Items = eventData.Items
	}
	if eventData.UpdateOptional() {
		milestone.Optional = eventData.Optional
	}
	if eventData.UpdateRetired() {
		milestone.Retired = eventData.Retired
	}

	a.MasterPlan.Milestones[milestone.ID] = milestone

	return nil
}

func (a *MasterPlanAggregate) onMasterPlanUpdate(evt eventstore.Event) error {
	var eventData event.MasterPlanUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if eventData.UpdateName() {
		a.MasterPlan.Name = eventData.Name
	}

	if eventData.UpdateRetired() {
		a.MasterPlan.Retired = eventData.Retired
	}

	return nil
}

func (a *MasterPlanAggregate) onMasterPlanMilestoneReorder(evt eventstore.Event) error {
	var eventData event.MasterPlanMilestoneReorderEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.MasterPlan.Milestones == nil {
		a.MasterPlan.Milestones = make(map[string]model.MasterPlanMilestone)
	}
	for i, milestoneId := range eventData.MilestoneIds {
		if milestone, ok := a.MasterPlan.Milestones[milestoneId]; ok {
			milestone.Order = int64(i)
			a.MasterPlan.Milestones[milestoneId] = milestone
		}
	}
	return nil
}
