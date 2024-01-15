package command

import (
	"time"

	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type CreateOrganizationPlanMilestoneCommand struct {
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

func NewCreateOrganizationPlanMilestoneCommand(organizationPlanId, tenant, loggedInUserId, milestoneId, name string, order, durationHours int64, items []string, optional bool, sourceFields commonmodel.Source, createdAt *time.Time) *CreateOrganizationPlanMilestoneCommand {
	return &CreateOrganizationPlanMilestoneCommand{
		BaseCommand:   eventstore.NewBaseCommand(organizationPlanId, tenant, loggedInUserId).WithAppSource(sourceFields.AppSource),
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
