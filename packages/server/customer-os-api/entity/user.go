package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type UserEntity struct {
	Id              string
	FirstName       string                 `neo4jDb:"property:firstName;lookupName:FIRST_NAME;supportCaseSensitive:true"`
	LastName        string                 `neo4jDb:"property:lastName;lookupName:LAST_NAME;supportCaseSensitive:true"`
	Name            string                 `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	CreatedAt       time.Time              `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt       time.Time              `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	Source          neo4jentity.DataSource `neo4jDb:"property:source;lookupName:SOURCE;supportCaseSensitive:false"`
	SourceOfTruth   neo4jentity.DataSource `neo4jDb:"property:sourceOfTruth;lookupName:SOURCE;supportCaseSensitive:false"`
	AppSource       string                 `neo4jDb:"property:appSource;lookupName:APP_SOURCE;supportCaseSensitive:false"`
	Roles           []string               `neo4jDb:"property:roles;lookupName:ROLES;supportCaseSensitive:false"`
	Timezone        string                 `neo4jDb:"property:timezone;lookupName:TIMEZONE;supportCaseSensitive:true"`
	ProfilePhotoUrl string                 `neo4jDb:"property:profilePhotoUrl;lookupName:PROFILE_PHOTO_URL;supportCaseSensitive:true"`
	Internal        bool
	Bot             bool

	DefaultForPlayer bool
	Tenant           string

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails
	DataloaderKey                        string
}

func (User UserEntity) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", User.Id, User.FirstName, User.LastName)
}

type UserEntities []UserEntity

func (UserEntity) IsInteractionEventParticipant() {}

func (UserEntity) IsInteractionSessionParticipant() {}

func (UserEntity) IsIssueParticipant() {}

func (UserEntity) IsMeetingParticipant() {}

func (UserEntity) ParticipantLabel() string {
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
