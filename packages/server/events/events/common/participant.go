package common

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
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
		return neo4jutil.NodeLabelUser
	case ContactType:
		return neo4jutil.NodeLabelContact
	case OrganizationType:
		return neo4jutil.NodeLabelOrganization
	case JobRoleType:
		return neo4jutil.NodeLabelJobRole
	default:
		return ""
	}
}
