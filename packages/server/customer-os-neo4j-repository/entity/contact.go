package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"strings"
	"time"
)

type ContactProperty string

const (
	ContactPropertyEnrichedAt                                ContactProperty = "enrichedAt"
	ContactPropertyEnrichedFailedAtScrapInPersonSearch       ContactProperty = "enrichedFailedAtScrapInPersonSearch"
	ContactPropertyEnrichedAtScrapInPersonSearch             ContactProperty = "enrichedAtScrapInPersonSearch"
	ContactPropertyEnrichedAtScrapInProfile                  ContactProperty = "enrichedAtScrapInProfile"
	ContactPropertyEnrichedScrapInPersonSearchParam          ContactProperty = "enrichedScrapInPersonSearchParam"
	ContactPropertyEnrichedScrapInProfileParam               ContactProperty = "enrichedScrapInProfileParam"
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
	BettercontactFoundEmailAt           *time.Time
	EnrichedAt                          *time.Time
	EnrichedAtScrapInPersonSearch       *time.Time
	EnrichedFailedAtScrapInPersonSearch *time.Time
	EnrichedScrapInPersonSearchParam    string
	EnrichedAtScrapInProfile            *time.Time
	EnrichedScrapInProfileParam         string
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

func (c ContactEntity) DeriveFirstAndLastNames() (string, string) {
	firstName := c.FirstName
	lastName := c.LastName
	if (firstName == "" || lastName == "") && c.Name != "" {
		parts := strings.Split(c.Name, " ")
		if firstName == "" {
			firstName = parts[0]
		}
		if lastName == "" && len(parts) > 1 {
			lastName = strings.Join(parts[1:], " ")
		}
	}
	return firstName, lastName
}
