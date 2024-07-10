package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"time"
)

type CreateMasterPlanMilestoneCommand struct {
	eventstore.BaseCommand
	SourceFields  events.Source
	MilestoneId   string
	CreatedAt     *time.Time
	Name          string
	Order         int64
	DurationHours int64
	Items         []string
	Optional      bool
}

func NewCreateMasterPlanMilestoneCommand(masterPlanId, tenant, loggedInUserId, milestoneId, name string, order, durationHours int64, items []string, optional bool, sourceFields events.Source, createdAt *time.Time) *CreateMasterPlanMilestoneCommand {
	return &CreateMasterPlanMilestoneCommand{
		BaseCommand:   eventstore.NewBaseCommand(masterPlanId, tenant, loggedInUserId).WithAppSource(sourceFields.AppSource),
		SourceFields:  sourceFields,
		CreatedAt:     createdAt,
		Name:          name,
		MilestoneId:   milestoneId,
		Order:         order,
		DurationHours: durationHours,
		Items:         items,
		Optional:      optional,
	}
}
