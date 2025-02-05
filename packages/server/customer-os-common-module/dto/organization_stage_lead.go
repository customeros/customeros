package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type OrganizationStageLead struct {
	Source         enum.Source
	OrganizationID string `json:"organizationId"`
}

func (e *OrganizationStageLead) Name() enum.AgentListenerEvent {
	return enum.EventCompanyStageLead
}
