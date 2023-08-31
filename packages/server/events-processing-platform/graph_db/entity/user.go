package entity

import (
	"fmt"
	"time"
)

type UserEntity struct {
	Id              string
	FirstName       string
	LastName        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Source          DataSource
	SourceOfTruth   DataSource
	AppSource       string
	Roles           []string
	ProfilePhotoUrl string
	Timezone        string
	Internal        bool

	DefaultForPlayer bool
	Tenant           string

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails
	DataloaderKey                        string
}

func (u UserEntity) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", u.Id, u.FirstName, u.LastName)
}

type UserEntities []UserEntity

func (UserEntity) IsInteractionEventParticipant() {}

func (UserEntity) InteractionEventParticipantLabel() string {
	return NodeLabel_User
}

func (UserEntity) IsInteractionSessionParticipant() {}

func (UserEntity) IsMeetingParticipant() {}

func (UserEntity) InteractionSessionParticipantLabel() string {
	return NodeLabel_User
}

func (UserEntity) MeetingParticipantLabel() string {
	return NodeLabel_User
}

func (user UserEntity) GetDataloaderKey() string {
	return user.DataloaderKey
}

func (UserEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_User,
		NodeLabel_User + "_" + tenant,
	}
}
