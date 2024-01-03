package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type ContactEntity struct {
	Id              string
	Prefix          string `neo4jDb:"property:prefix;lookupName:PREFIX;supportCaseSensitive:true"`
	Name            string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	FirstName       string `neo4jDb:"property:firstName;lookupName:FIRST_NAME;supportCaseSensitive:true"`
	LastName        string `neo4jDb:"property:lastName;lookupName:LAST_NAME;supportCaseSensitive:true"`
	Description     string `neo4jDb:"property:description;lookupName:DESCRIPTION;supportCaseSensitive:true"`
	Timezone        string `neo4jDb:"property:timezone;lookupName:TIMEZONE;supportCaseSensitive:true"`
	ProfilePhotoUrl string `neo4jDb:"property:profilePhotoUrl;lookupName:PROFILE_PHOTO_URL;supportCaseSensitive:true"`
	// TODO make non pointer
	CreatedAt     *time.Time             `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt     time.Time              `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	Source        neo4jentity.DataSource `neo4jDb:"property:source;lookupName:SOURCE;supportCaseSensitive:false"`
	SourceOfTruth neo4jentity.DataSource `neo4jDb:"property:sourceOfTruth;lookupName:SOURCE_OF_TRUTH;supportCaseSensitive:false"`
	AppSource     string                 `neo4jDb:"property:appSource;lookupName:APP_SOURCE;supportCaseSensitive:false"`

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails
	DataloaderKey                        string
}

func (contact ContactEntity) ToString() string {
	return fmt.Sprintf("ContactEntity{Id: %s, Prefix: %s, Name: %s, FirstName: %s, LastName: %s, Description: %s, CreatedAt: %s, UpdatedAt: %s, Source: %s, SourceOfTruth: %s, AppSource: %s}",
		contact.Id, contact.Prefix, contact.Name, contact.FirstName, contact.LastName, contact.Description, contact.CreatedAt, contact.UpdatedAt, contact.Source, contact.SourceOfTruth, contact.AppSource)
}

type ContactEntities []ContactEntity

func (ContactEntity) IsInteractionEventParticipant() {}

func (ContactEntity) IsInteractionSessionParticipant() {}

func (ContactEntity) IsIssueParticipant() {}

func (ContactEntity) ParticipantLabel() string {
	return NodeLabel_Contact
}

func (ContactEntity) IsMeetingParticipant() {}

func (ContactEntity) MeetingParticipantLabel() string {
	return NodeLabel_Contact
}

func (ContactEntity) IsNotedEntity() {}

func (ContactEntity) NotedEntityLabel() string {
	return NodeLabel_Contact
}

func (contact ContactEntity) GetDataloaderKey() string {
	return contact.DataloaderKey
}

func (ContactEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Contact,
		NodeLabel_Contact + "_" + tenant,
	}
}
