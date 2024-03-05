package command

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CreateMasterPlanMilestoneCommand struct {
	eventstore.BaseCommand
	SourceFields  commonmodel.Source
	MilestoneId   string
	CreatedAt     *time.Time
	Name          string
	Order         int64
	DurationHours int64
	Items         []string
	Optional      bool
}

func NewCreateMasterPlanMilestoneCommand(masterPlanId, tenant, loggedInUserId, milestoneId, name string, order, durationHours int64, items []string, optional bool, sourceFields commonmodel.Source, createdAt *time.Time) *CreateMasterPlanMilestoneCommand {
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
