package command

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CreateRenewalOpportunityCommand struct {
	eventstore.BaseCommand
	Source     commonmodel.Source
	ContractId string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func NewCreateRenewalOpportunityCommand(opportunityId, tenant, loggedInUserId, contractId string, source commonmodel.Source, createdAt, updatedAt *time.Time) *CreateRenewalOpportunityCommand {
	return &CreateRenewalOpportunityCommand{
		BaseCommand: eventstore.NewBaseCommand(opportunityId, tenant, loggedInUserId),
		Source:      source,
		ContractId:  contractId,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
