package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	MasterPlanAggregateType eventstore.AggregateType = "master_plan"
)

type MasterPlanAggregate struct {
	*aggregate.CommonTenantIdAggregate
	MasterPlan *model.MasterPlan
}

func NewMasterPlanAggregateWithTenantAndID(tenant, id string) *MasterPlanAggregate {
	masterPlanAggregate := MasterPlanAggregate{}
	masterPlanAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(MasterPlanAggregateType, tenant, id)
	masterPlanAggregate.SetWhen(masterPlanAggregate.When)
	masterPlanAggregate.MasterPlan = &model.MasterPlan{}
	masterPlanAggregate.Tenant = tenant

	return &masterPlanAggregate
}

func (a *MasterPlanAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.MasterPlanCreateV1:
		return a.onMasterPlanCreate(evt)
	default:
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
