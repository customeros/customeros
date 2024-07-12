package common

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
)

type ParticipantType string

const (
	UserType         ParticipantType = "User"
	ContactType      ParticipantType = "Contact"
	OrganizationType ParticipantType = "Organization"
	JobRoleType      ParticipantType = "JobRole"
)

type Participant struct {
	ID              string          `json:"id"`
	ParticipantType ParticipantType `json:"participantType"`
}

func (p Participant) NodeLabel() string {
	switch p.ParticipantType {
	case UserType:
		return model.NodeLabelUser
	case ContactType:
		return model.NodeLabelContact
	case OrganizationType:
		return model.NodeLabelOrganization
	case JobRoleType:
		return model.NodeLabelJobRole
	default:
		return ""
	}
}
