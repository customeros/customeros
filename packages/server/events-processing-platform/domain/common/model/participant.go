package model

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
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
		return neo4jentity.NodeLabel_User
	case ContactType:
		return neo4jentity.NodeLabel_Contact
	case OrganizationType:
		return neo4jentity.NodeLabel_Organization
	default:
		return ""
	}
}
