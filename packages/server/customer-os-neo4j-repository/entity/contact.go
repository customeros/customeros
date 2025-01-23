package neo4j_entity

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
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

func (c ContactEntity) IsHidden() bool {
	return c.Hide
}

func (c ContactEntity) FullName() string {
	first := strings.TrimSpace(c.FirstName)
	last := strings.TrimSpace(c.LastName)

	if first != "" && last != "" {
		return first + " " + last
	} else if first != "" {
		return first
	} else if last != "" {
		return last
	}
	return ""
}
