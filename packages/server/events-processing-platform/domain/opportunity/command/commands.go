package command

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

// CreateOpportunityCommand contains the data needed to create an opportunity.
type CreateOpportunityCommand struct {
	eventstore.BaseCommand
	DataFields     model.OpportunityDataFields
	Source         commonmodel.Source
	ExternalSystem commonmodel.ExternalSystem
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

// NewCreateOpportunityCommand creates a new CreateOpportunityCommand.
func NewCreateOpportunityCommand(opportunityId, tenant, loggedInUserId string, dataFields model.OpportunityDataFields, source commonmodel.Source, externalSystem commonmodel.ExternalSystem, createdAt, updatedAt *time.Time) *CreateOpportunityCommand {
	return &CreateOpportunityCommand{
		BaseCommand:    eventstore.NewBaseCommand(opportunityId, tenant, loggedInUserId),
		DataFields:     dataFields,
		Source:         source,
		ExternalSystem: externalSystem,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

// UpdateOpportunityCommand contains the data needed to update an opportunity.
type UpdateOpportunityCommand struct {
	eventstore.BaseCommand
	DataFields     model.OpportunityDataFields
	ExternalSystem commonmodel.ExternalSystem
	Source         commonmodel.Source
	UpdatedAt      *time.Time
}

// NewUpdateOpportunityCommand creates a new UpdateOpportunityCommand.
func NewUpdateOpportunityCommand(opportunityId, tenant, loggedInUserId string, dataFields model.OpportunityDataFields, source commonmodel.Source, externalSystem commonmodel.ExternalSystem, updatedAt *time.Time) *UpdateOpportunityCommand {
	return &UpdateOpportunityCommand{
		BaseCommand:    eventstore.NewBaseCommand(opportunityId, tenant, loggedInUserId),
		DataFields:     dataFields,
		ExternalSystem: externalSystem,
		UpdatedAt:      updatedAt,
		Source:         source,
	}
}
