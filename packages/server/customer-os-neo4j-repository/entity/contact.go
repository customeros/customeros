package entity

import (
	"github.com/forPelevin/gomoji"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
	"time"
)

type ContactProperty string

const (
	ContactPropertyEnrichedAt                                  ContactProperty = "enrichedAt"
	ContactPropertyEnrichFailedAt                              ContactProperty = "enrichFailedAt"
	ContactPropertyEnrichRequestedAt                           ContactProperty = "techEnrichRequestedAt"
	ContactPropertyEnrichAttempts                              ContactProperty = "techEnrichAttempts"
	ContactPropertyEnrichedScrapinRecordId                     ContactProperty = "enrichedScrapinRecordId"
	ContactPropertyBettercontactFoundEmailAt                   ContactProperty = "bettercontactFoundEmailAt"
	ContactPropertyFindWorkEmailWithBetterContactRequestedId   ContactProperty = "techFindWorkEmailWithBetterContactRequestId"
	ContactPropertyFindWorkEmailWithBetterContactRequestedAt   ContactProperty = "techFindWorkEmailWithBetterContactRequestedAt"
	ContactPropertyFindWorkEmailWithBetterContactCompletedAt   ContactProperty = "techFindWorkEmailWithBetterContactCompletedAt"
	ContactPropertyFindWorkEmailWithBetterContactFound         ContactProperty = "techFindWorkEmailWithBetterContactFound"
	ContactPropertyFindMobilePhoneWithBetterContactRequestedId ContactProperty = "techFindMobilePhoneWithBetterContactRequestId"
	ContactPropertyFindMobilePhoneWithBetterContactRequestedAt ContactProperty = "techFindMobilePhoneWithBetterContactRequestedAt"
	ContactPropertyFindMobilePhoneWithBetterContactCompletedAt ContactProperty = "techFindMobilePhoneWithBetterContactCompletedAt"
	ContactPropertyFindMobilePhoneWithBetterContactFound       ContactProperty = "techFindMobilePhoneWithBetterContactFound"
	ContactPropertyUpdateWithWorkEmailRequestedAt              ContactProperty = "techUpdateWithWorkEmailRequestedAt"
	ContactPropertyCheckedAt                                   ContactProperty = "techCheckedAt"
	ContactPropertyPrefix                                      ContactProperty = "prefix"
	ContactPropertyName                                        ContactProperty = "name"
	ContactPropertyFirstName                                   ContactProperty = "firstName"
	ContactPropertyLastName                                    ContactProperty = "lastName"
	ContactPropertyDescription                                 ContactProperty = "description"
	ContactPropertyTimezone                                    ContactProperty = "timezone"
	ContactPropertyProfilePhotoUrl                             ContactProperty = "profilePhotoUrl"
	ContactPropertyHide                                        ContactProperty = "hide"
	ContactPropertyUsername                                    ContactProperty = "username"
	ContactPropertyHiddenAt                                    ContactProperty = "hiddenAt"
)

type ContactEntity struct {
	DataLoaderKey
	EventStoreAggregate
	Id        string
	CreatedAt time.Time `neo4jDb:"property:createdAt;lookupName:CREATED_AT"`
	UpdatedAt time.Time `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT"`
	Source    DataSource
	AppSource string

	Prefix          string `neo4jDb:"property:prefix;lookupName:PREFIX;supportCaseSensitive:true"`
	Name            string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	FirstName       string `neo4jDb:"property:firstName;lookupName:FIRST_NAME;supportCaseSensitive:true"`
	LastName        string `neo4jDb:"property:lastName;lookupName:LAST_NAME;supportCaseSensitive:true"`
	Description     string `neo4jDb:"property:description;lookupName:DESCRIPTION;supportCaseSensitive:true"`
	Timezone        string `neo4jDb:"property:timezone;lookupName:TIMEZONE;supportCaseSensitive:true"`
	ProfilePhotoUrl string `neo4jDb:"property:profilePhotoUrl;lookupName:PROFILE_PHOTO_URL;supportCaseSensitive:true"`
	Hide            bool   `neo4jDb:"property:hide;lookupName:HIDE;supportCaseSensitive:false"`
	Username        string `neo4jDb:"property:username;lookupName:USERNAME;supportCaseSensitive:true"`
	HiddenAt        *time.Time

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails

	ContactInternalFields ContactInternalFields
	EnrichDetails         ContactEnrichDetails
}

type ContactInternalFields struct {
	UpdateWithWorkEmailRequestedAt *time.Time
	HiddenAt                       *time.Time
	CheckedAt                      *time.Time
}

type ContactEnrichDetails struct {
	EnrichRequestedAt                           *time.Time
	EnrichedAt                                  *time.Time
	EnrichFailedAt                              *time.Time
	EnrichAttempts                              int64
	BettercontactFoundEmailAt                   *time.Time // TODO alexb check how it is used
	EnrichedScrapinRecordId                     string
	FindWorkEmailWithBetterContactRequestedId   *string
	FindWorkEmailWithBetterContactRequestedAt   *time.Time
	FindWorkEmailWithBetterContactCompletedAt   *time.Time
	FindWorkEmailWithBetterContactFound         *bool
	FindMobilePhoneWithBetterContactRequestedId *string
	FindMobilePhoneWithBetterContactRequestedAt *time.Time
	FindMobilePhoneWithBetterContactCompletedAt *time.Time
	FindMobilePhoneWithBetterContactFound       *bool
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

func (c ContactEntity) GetNamesFromString(input string) (string, string) {
	firstName := ""
	lastName := ""
	specialChars := []string{" ", ".", "-", "_", "+", "=", ","}

	// Trim spaces
	input = gomoji.RemoveEmojis(input)
	input = strings.TrimSpace(input)

	if input != "" {
		// Find the position of the first occurrence of any special character
		splitPos := -1
		for _, char := range specialChars {
			if pos := strings.Index(input, char); pos != -1 {
				if splitPos == -1 || pos < splitPos {
					splitPos = pos
				}
			}
		}

		if splitPos == -1 {
			// No special characters, treat input as a single word (set as first name)
			firstName = input
		} else {
			// Split input at the first special character
			firstName = input[:splitPos]
			lastName = strings.TrimSpace(input[splitPos+1:])
		}

		// Apply Camel case to both names
		firstName = utils.ToCamelCase(firstName)
		lastName = utils.ToCamelCase(lastName)
	}

	return firstName, lastName
}
