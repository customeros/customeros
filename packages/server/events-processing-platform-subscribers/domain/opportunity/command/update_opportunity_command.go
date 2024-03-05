package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpdateOpportunityCommand struct {
	eventstore.BaseCommand
	DataFields     model.OpportunityDataFields
	ExternalSystem commonmodel.ExternalSystem
	Source         commonmodel.Source
	UpdatedAt      *time.Time
	FieldsMask     []string
}

func NewUpdateOpportunityCommand(opportunityId, tenant, loggedInUserId string, dataFields model.OpportunityDataFields, source commonmodel.Source, externalSystem commonmodel.ExternalSystem, updatedAt *time.Time, fieldsMask []string) *UpdateOpportunityCommand {
	return &UpdateOpportunityCommand{
		BaseCommand:    eventstore.NewBaseCommand(opportunityId, tenant, loggedInUserId),
		ExternalSystem: externalSystem,
		UpdatedAt:      updatedAt,
		Source:         source,
		DataFields:     dataFields,
		FieldsMask:     fieldsMask,
	}
}

func (c *UpdateOpportunityCommand) FieldMaskContainsOnly(fieldMasks []string) bool {
	if len(c.FieldsMask) != len(fieldMasks) {
		return false
	}
	for _, fieldMask := range fieldMasks {
		if !utils.Contains(c.FieldsMask, fieldMask) {
			return false
		}
	}
	return true
}
