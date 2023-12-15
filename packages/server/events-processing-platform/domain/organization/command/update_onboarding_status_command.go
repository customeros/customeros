package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"strings"
	"time"
)

type UpdateOnboardingStatusCommand struct {
	eventstore.BaseCommand
	Status             string     `json:"status"`
	Comments           string     `json:"comments"`
	UpdatedAt          *time.Time `json:"updatedAt"`
	CausedByContractId string     `json:"causedByContractId"`
}

func NewUpdateOnboardingStatusCommand(tenant, orgId, loggedInUserId, appSource, status, comments, causedByContractId string, updatedAt *time.Time) *UpdateOnboardingStatusCommand {
	return &UpdateOnboardingStatusCommand{
		BaseCommand:        eventstore.NewBaseCommand(orgId, tenant, loggedInUserId).WithAppSource(appSource),
		Status:             status,
		Comments:           strings.TrimSpace(comments),
		UpdatedAt:          updatedAt,
		CausedByContractId: strings.TrimSpace(causedByContractId),
	}
}
