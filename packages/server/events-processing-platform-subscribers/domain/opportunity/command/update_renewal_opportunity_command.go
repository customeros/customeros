package command

import (
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpdateRenewalOpportunityCommand struct {
	eventstore.BaseCommand
	Source            commonmodel.Source
	UpdatedAt         *time.Time
	RenewalLikelihood neo4jenum.RenewalLikelihood
	Comments          string
	Amount            float64
	MaskFields        []string
	OwnerUserId       string
}

func NewUpdateRenewalOpportunityCommand(opportunityId, tenant, loggedInUserId, comments string, renewalLikelihood neo4jenum.RenewalLikelihood, amount float64, source commonmodel.Source, updatedAt *time.Time, maskFields []string, ownerUserId string) *UpdateRenewalOpportunityCommand {
	return &UpdateRenewalOpportunityCommand{
		BaseCommand:       eventstore.NewBaseCommand(opportunityId, tenant, loggedInUserId),
		Source:            source,
		RenewalLikelihood: renewalLikelihood,
		UpdatedAt:         updatedAt,
		Comments:          comments,
		Amount:            amount,
		MaskFields:        maskFields,
		OwnerUserId:       ownerUserId,
	}
}
