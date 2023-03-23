package entity

import (
	"fmt"
	"time"
)

type ContactEntity struct {
	Id        string
	Title     string `neo4jDb:"property:title;lookupName:TITLE;supportCaseSensitive:false"`
	Name      string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	FirstName string `neo4jDb:"property:firstName;lookupName:FIRST_NAME;supportCaseSensitive:true"`
	LastName  string `neo4jDb:"property:lastName;lookupName:LAST_NAME;supportCaseSensitive:true"`
	// TODO make non pointer and use different property for input
	CreatedAt     *time.Time `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt     time.Time  `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	Source        DataSource `neo4jDb:"property:source;lookupName:SOURCE;supportCaseSensitive:false"`
	SourceOfTruth DataSource `neo4jDb:"property:sourceOfTruth;lookupName:SOURCE_OF_TRUTH;supportCaseSensitive:false"`
	AppSource     string     `neo4jDb:"property:appSource;lookupName:APP_SOURCE;supportCaseSensitive:false"`

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails
	DataloaderKey                        string
}

func (contact ContactEntity) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", contact.Id, contact.FirstName, contact.LastName)
}

type ContactEntities []ContactEntity

func (ContactEntity) IsInteractionEventParticipant() {}

func (ContactEntity) InteractionEventParticipantLabel() string {
	return NodeLabel_Contact
}

func (ContactEntity) IsInteractionSessionParticipant() {}

func (ContactEntity) InteractionSessionParticipantLabel() string {
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
