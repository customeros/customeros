package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"strings"
	"time"
)

type ContactProperty string

const (
	ContactPropertyEnrichedAt                                ContactProperty = "enrichedAt"
	ContactPropertyEnrichFailedAt                            ContactProperty = "enrichFailedAt"
	ContactPropertyBettercontactFoundEmailAt                 ContactProperty = "bettercontactFoundEmailAt"
	ContactPropertyFindWorkEmailWithBetterContactRequestedId ContactProperty = "techFindWorkEmailWithBetterContactRequestId"
	ContactPropertyFindWorkEmailWithBetterContactRequestedAt ContactProperty = "techFindWorkEmailWithBetterContactRequestedAt"
	ContactPropertyFindWorkEmailWithBetterContactCompletedAt ContactProperty = "techFindWorkEmailWithBetterContactCompletedAt"
	ContactPropertyEnrichRequestedAt                         ContactProperty = "techEnrichRequestedAt"
	ContactPropertyPrefix                                    ContactProperty = "prefix"
	ContactPropertyName                                      ContactProperty = "name"
	ContactPropertyFirstName                                 ContactProperty = "firstName"
	ContactPropertyLastName                                  ContactProperty = "lastName"
	ContactPropertyDescription                               ContactProperty = "description"
	ContactPropertyTimezone                                  ContactProperty = "timezone"
	ContactPropertyProfilePhotoUrl                           ContactProperty = "profilePhotoUrl"
	ContactPropertyHide                                      ContactProperty = "hide"
	ContactPropertyUsername                                  ContactProperty = "username"
	ContactPropertyEnrichedScrapinRecordId                   ContactProperty = "enrichedScrapinRecordId"
)

type ContactEntity struct {
	DataLoaderKey
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
	Hide            bool   `neo4jDb:"property:hide;lookupName:HIDE;supportCaseSensitive:false"`
	Username        string `neo4jDb:"property:username;lookupName:USERNAME;supportCaseSensitive:true"`

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails

	ContactInternalFields ContactInternalFields
	EnrichDetails         ContactEnrichDetails
}

type ContactInternalFields struct {
	FindWorkEmailWithBetterContactRequestedId *string
	FindWorkEmailWithBetterContactRequestedAt *time.Time
	FindWorkEmailWithBetterContactCompletedAt *time.Time
	EnrichRequestedAt                         *time.Time
}

type ContactEnrichDetails struct {
	BettercontactFoundEmailAt *time.Time
	EnrichedAt                *time.Time
	EnrichedFailedAt          *time.Time
	EnrichedScrapinRecordId   string
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
	return model.NodeLabelContact
}

func (c ContactEntity) Labels(tenant string) []string {
	return []string{c.EntityLabel(), c.EntityLabel() + "_" + tenant}
}

func (c ContactEntity) DeriveFirstAndLastNames() (string, string) {
	firstName := strings.TrimSpace(c.FirstName)
	lastName := strings.TrimSpace(c.LastName)
	name := strings.TrimSpace(c.Name)
	if (firstName == "" || lastName == "") && name != "" {
		parts := strings.Split(name, " ")
		if firstName == "" {
			firstName = parts[0]
		}
		if lastName == "" && len(parts) > 1 {
			lastName = strings.Join(parts[1:], " ")
		}
	}

	if firstName != "" && lastName == "" {
		parts := strings.Split(firstName, " ")
		if len(parts) > 1 {
			firstName = parts[0]
			lastName = strings.Join(parts[1:], " ")
		}
	}

	if firstName == "" && lastName != "" {
		parts := strings.Split(lastName, " ")
		if len(parts) > 1 {
			firstName = parts[0]
			lastName = strings.Join(parts[1:], " ")
		}
	}

	return firstName, lastName
}
