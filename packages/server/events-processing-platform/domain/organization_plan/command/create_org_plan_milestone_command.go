package command

import (
	"time"

	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type CreateOrgPlanMilestoneCommand struct {
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

func NewCreateOrgPlanMilestoneCommand(orgPlanId, tenant, loggedInUserId, milestoneId, name string, order, durationHours int64, items []string, optional bool, sourceFields commonmodel.Source, createdAt *time.Time) *CreateOrgPlanMilestoneCommand {
	return &CreateOrgPlanMilestoneCommand{
		BaseCommand:   eventstore.NewBaseCommand(orgPlanId, tenant, loggedInUserId).WithAppSource(sourceFields.AppSource),
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
