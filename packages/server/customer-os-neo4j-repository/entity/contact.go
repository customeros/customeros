package entity

import (
	"time"
)

type ContactEntity struct {
	EventStoreAggregate
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	Prefix    string `neo4jDb:"property:prefix;lookupName:PREFIX;supportCaseSensitive:true"`
	Name      string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	FirstName string `neo4jDb:"property:firstName;lookupName:FIRST_NAME;supportCaseSensitive:true"`
	LastName  string `neo4jDb:"property:lastName;lookupName:LAST_NAME;supportCaseSensitive:true"`

	Description     string `neo4jDb:"property:description;lookupName:DESCRIPTION;supportCaseSensitive:true"`
	Timezone        string `neo4jDb:"property:timezone;lookupName:TIMEZONE;supportCaseSensitive:true"`
	ProfilePhotoUrl string `neo4jDb:"property:profilePhotoUrl;lookupName:PROFILE_PHOTO_URL;supportCaseSensitive:true"`

	DataloaderKey string
}

type ContactEntities []ContactEntity
