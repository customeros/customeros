package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
)

type ParticipantType string

const (
	UserType         ParticipantType = "User"
	ContactType      ParticipantType = "Contact"
	OrganizationType ParticipantType = "Organization"
)

type Participant struct {
	ID              string          `json:"id"`
	ParticipantType ParticipantType `json:"participantType"`
}

func (p Participant) NodeLabel() string {
	switch p.ParticipantType {
	case UserType:
		return constants.NodeLabel_User
	case ContactType:
		return constants.NodeLabel_Contact
	case OrganizationType:
		return constants.NodeLabel_Organization
	default:
		return ""
	}
}
