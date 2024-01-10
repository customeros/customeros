package command

import (
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CreateRenewalOpportunityCommand struct {
	eventstore.BaseCommand
	Source            commonmodel.Source
	ContractId        string
	RenewalLikelihood neo4jenum.RenewalLikelihood
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
}

func NewCreateRenewalOpportunityCommand(opportunityId, tenant, loggedInUserId, contractId string, renewalLikelihood neo4jenum.RenewalLikelihood, source commonmodel.Source, createdAt, updatedAt *time.Time) *CreateRenewalOpportunityCommand {
	return &CreateRenewalOpportunityCommand{
		BaseCommand:       eventstore.NewBaseCommand(opportunityId, tenant, loggedInUserId),
		Source:            source,
		ContractId:        contractId,
		RenewalLikelihood: renewalLikelihood,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}
