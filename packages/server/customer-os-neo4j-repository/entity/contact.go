package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type ContactEntity struct {
	EventStoreAggregate
	Id            string
	CreatedAt     time.Time `neo4jDb:"property:createdAt;lookupName:CREATED_AT"`
	UpdatedAt     time.Time `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT"`
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	Prefix          string `neo4jDb:"property:prefix;lookupName:PREFIX;supportCaseSensitive:true"`
	Name            string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	FirstName       string `neo4jDb:"property:firstName;lookupName:FIRST_NAME;supportCaseSensitive:true"`
	LastName        string `neo4jDb:"property:lastName;lookupName:LAST_NAME;supportCaseSensitive:true"`
	Description     string `neo4jDb:"property:description;lookupName:DESCRIPTION;supportCaseSensitive:true"`
	Timezone        string `neo4jDb:"property:timezone;lookupName:TIMEZONE;supportCaseSensitive:true"`
	ProfilePhotoUrl string `neo4jDb:"property:profilePhotoUrl;lookupName:PROFILE_PHOTO_URL;supportCaseSensitive:true"`

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails

	DataloaderKey string
}

type ContactEntities []ContactEntity

func (c ContactEntity) GetDataloaderKey() string {
	return c.DataloaderKey
}

func (ContactEntity) IsIssueParticipant() {}

func (ContactEntity) IsInteractionEventParticipant() {}

func (ContactEntity) IsInteractionSessionParticipant() {}

func (ContactEntity) IsMeetingParticipant() {}

func (ContactEntity) EntityLabel() string {
	return neo4jutil.NodeLabelContact
}

func (c ContactEntity) Labels(tenant string) []string {
	return []string{c.EntityLabel(), c.EntityLabel() + "_" + tenant}
}
