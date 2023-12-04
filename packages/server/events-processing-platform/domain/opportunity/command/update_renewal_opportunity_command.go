package command

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpdateRenewalOpportunityCommand struct {
	eventstore.BaseCommand
	Source            commonmodel.Source
	UpdatedAt         *time.Time
	RenewalLikelihood model.RenewalLikelihoodString
	Comments          string
	Amount            float64
	MaskFields        []string
}

func NewUpdateRenewalOpportunityCommand(opportunityId, tenant, loggedInUserId, comments string, renewalLikelihood model.RenewalLikelihoodString, amount float64, source commonmodel.Source, updatedAt *time.Time, maskFields []string) *UpdateRenewalOpportunityCommand {
	return &UpdateRenewalOpportunityCommand{
		BaseCommand:       eventstore.NewBaseCommand(opportunityId, tenant, loggedInUserId),
		Source:            source,
		RenewalLikelihood: renewalLikelihood,
		UpdatedAt:         updatedAt,
		Comments:          comments,
		Amount:            amount,
		MaskFields:        maskFields,
	}
}
