package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/grpc/events_platform"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	eventcompletionpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_completion"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
	"time"
)

func TestQueryResolver_ContactByEmail(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	otherTenant := "other"
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenant(ctx, driver, otherTenant)
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, otherTenant)
	neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId1, "test@test.com", true, "MAIN")
	neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, otherTenant, contactId2, "test@test.com", true, "MAIN")

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_by_email"), client.Var("email", "test@test.com"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact_ByEmail model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, contactId1, contact.Contact_ByEmail.ID)
}

func TestQueryResolver_ContactByPhone(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	otherTenant := "other"
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenant(ctx, driver, otherTenant)
	contactId1 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, otherTenant)
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId1, "+1234567890", false, "OTHER")
	neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId2, "+1234567890", true, "MAIN")

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_by_phone"), client.Var("e164", "+1234567890"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contact struct {
		Contact_ByPhone model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contact)
	require.Nil(t, err)
	require.NotNil(t, contact)
	require.Equal(t, contactId1, contact.Contact_ByPhone.ID)
}

func TestQueryResolver_Contact_WithJobRoles_ById(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId1 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Name:        "name1",
		Description: "description1",
		Website:     "website1",
		Industry:    "industry1",
		IsPublic:    true,
	})
	neo4jt.AddDomainToOrg(ctx, driver, organizationId1, "domain1")
	organizationId2 := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Name:        "name2",
		Description: "description2",
		Website:     "website2",
		Industry:    "industry2",
		IsPublic:    false,
	})
	neo4jt.AddDomainToOrg(ctx, driver, organizationId2, "domain2")
	role1 := neo4jt.ContactWorksForOrganization(ctx, driver, contactId, organizationId1, "CTO", false)
	role2 := neo4jt.ContactWorksForOrganization(ctx, driver, contactId, organizationId2, "CEO", true)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "JobRole"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "WORKS_AS"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_job_roles_by_id"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var searchedContact struct {
		Contact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &searchedContact)
	require.Nil(t, err)
	require.Equal(t, contactId, searchedContact.Contact.ID)

	roles := searchedContact.Contact.JobRoles
	require.Equal(t, 2, len(roles))
	var cto, ceo *model.JobRole
	ceo = roles[0]
	cto = roles[1]
	require.Equal(t, role1, cto.ID)
	require.Equal(t, "CTO", *cto.JobTitle)
	require.Equal(t, false, cto.Primary)
	require.Equal(t, organizationId1, cto.Organization.ID)
	require.Equal(t, "name1", cto.Organization.Name)
	require.Equal(t, "description1", *cto.Organization.Description)
	require.Equal(t, []string{"domain1"}, cto.Organization.Domains)
	require.Equal(t, "website1", *cto.Organization.Website)
	require.Equal(t, "industry1", *cto.Organization.Industry)
	require.Equal(t, true, *cto.Organization.IsPublic)
	require.NotNil(t, cto.Organization.CreatedAt)

	require.Equal(t, role2, ceo.ID)
	require.Equal(t, "CEO", *ceo.JobTitle)
	require.Equal(t, true, ceo.Primary)
	require.Equal(t, organizationId2, ceo.Organization.ID)
	require.Equal(t, "name2", ceo.Organization.Name)
	require.Equal(t, "description2", *ceo.Organization.Description)
	require.Equal(t, []string{"domain2"}, ceo.Organization.Domains)
	require.Equal(t, "website2", *ceo.Organization.Website)
	require.Equal(t, "industry2", *ceo.Organization.Industry)
	require.Equal(t, false, *ceo.Organization.IsPublic)
	require.NotNil(t, ceo.Organization.CreatedAt)
}

func TestQueryResolver_Contact_WithTags_ById(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	tagId1 := neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{
		Name:      "tag1",
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
		Source:    neo4jentity.DataSourceOpenline,
	})
	tagId2 := neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{
		Name:      "tag2",
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
		Source:    neo4jentity.DataSourceOpenline,
	})
	tagId3 := neo4jtest.CreateTag(ctx, driver, tenantName, neo4jentity.TagEntity{
		Name:      "tag3",
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
		Source:    neo4jentity.DataSourceOpenline,
	})
	neo4jt.TagContact(ctx, driver, contactId, tagId1)
	neo4jt.TagContact(ctx, driver, contactId, tagId2)
	neo4jt.TagContact(ctx, driver, contactId2, tagId3)

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Tag"))
	require.Equal(t, 3, neo4jtest.GetCountOfRelationships(ctx, driver, "TAGGED"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_tags_by_id"),
		client.Var("contactId", contactId))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactStruct struct {
		Contact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	contact := contactStruct.Contact

	require.Nil(t, err)
	require.Equal(t, contactId, contact.ID)

	tags := contact.Tags
	require.Equal(t, 2, len(tags))
	require.Equal(t, tagId1, tags[0].ID)
	require.Equal(t, "tag1", tags[0].Name)
	require.Equal(t, tagId2, tags[1].ID)
	require.Equal(t, "tag2", tags[1].Name)
}

func TestQueryResolver_Contact_WithLocations_ById(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	locationId1 := neo4jt.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{
		Name:         "WORK",
		Source:       neo4jentity.DataSourceOpenline,
		AppSource:    "test",
		Country:      "testCountry",
		Region:       "testRegion",
		Locality:     "testLocality",
		Address:      "testAddress",
		Address2:     "testAddress2",
		Zip:          "testZip",
		AddressType:  "testAddressType",
		HouseNumber:  "testHouseNumber",
		PostalCode:   "testPostalCode",
		PlusFour:     "testPlusFour",
		Commercial:   true,
		Predirection: "testPredirection",
		District:     "testDistrict",
		Street:       "testStreet",
		RawAddress:   "testRawAddress",
		UtcOffset:    utils.Float64Ptr(1.0),
		TimeZone:     "paris",
		Latitude:     utils.ToPtr(float64(0.001)),
		Longitude:    utils.ToPtr(float64(-2.002)),
	})
	locationId2 := neo4jt.CreateLocation(ctx, driver, tenantName, neo4jentity.LocationEntity{
		Name:      "UNKNOWN",
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: "test",
	})
	neo4jt.ContactAssociatedWithLocation(ctx, driver, contactId, locationId1)
	neo4jt.ContactAssociatedWithLocation(ctx, driver, contactId, locationId2)

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_locations_by_id"),
		client.Var("contactId", contactId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var contactStruct struct {
		Contact model.Contact
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)

	contact := contactStruct.Contact
	require.NotNil(t, contact)
	require.Equal(t, 2, len(contact.Locations))

	var locationWithAddressDtls, locationWithoutAddressDtls *model.Location
	if contact.Locations[0].ID == locationId1 {
		locationWithAddressDtls = contact.Locations[0]
		locationWithoutAddressDtls = contact.Locations[1]
	} else {
		locationWithAddressDtls = contact.Locations[1]
		locationWithoutAddressDtls = contact.Locations[0]
	}

	require.Equal(t, locationId1, locationWithAddressDtls.ID)
	require.Equal(t, "WORK", *locationWithAddressDtls.Name)
	require.NotNil(t, locationWithAddressDtls.CreatedAt)
	require.NotNil(t, locationWithAddressDtls.UpdatedAt)
	require.Equal(t, "test", locationWithAddressDtls.AppSource)
	require.Equal(t, model.DataSourceOpenline, locationWithAddressDtls.Source)
	require.Equal(t, "testCountry", *locationWithAddressDtls.Country)
	require.Equal(t, "testLocality", *locationWithAddressDtls.Locality)
	require.Equal(t, "testRegion", *locationWithAddressDtls.Region)
	require.Equal(t, "testAddress", *locationWithAddressDtls.Address)
	require.Equal(t, "testAddress2", *locationWithAddressDtls.Address2)
	require.Equal(t, "testZip", *locationWithAddressDtls.Zip)
	require.Equal(t, "testAddressType", *locationWithAddressDtls.AddressType)
	require.Equal(t, "testHouseNumber", *locationWithAddressDtls.HouseNumber)
	require.Equal(t, "testPostalCode", *locationWithAddressDtls.PostalCode)
	require.Equal(t, "testPlusFour", *locationWithAddressDtls.PlusFour)
	require.Equal(t, true, *locationWithAddressDtls.Commercial)
	require.Equal(t, "testPredirection", *locationWithAddressDtls.Predirection)
	require.Equal(t, "testDistrict", *locationWithAddressDtls.District)
	require.Equal(t, "testStreet", *locationWithAddressDtls.Street)
	require.Equal(t, "testRawAddress", *locationWithAddressDtls.RawAddress)
	require.Equal(t, "paris", *locationWithAddressDtls.TimeZone)
	require.Equal(t, float64(1), *locationWithAddressDtls.UtcOffset)
	require.Equal(t, float64(0.001), *locationWithAddressDtls.Latitude)
	require.Equal(t, float64(-2.002), *locationWithAddressDtls.Longitude)

	require.Equal(t, locationId2, locationWithoutAddressDtls.ID)
	require.Equal(t, "UNKNOWN", *locationWithoutAddressDtls.Name)
	require.NotNil(t, locationWithoutAddressDtls.CreatedAt)
	require.NotNil(t, locationWithoutAddressDtls.UpdatedAt)
	require.Equal(t, "test", locationWithoutAddressDtls.AppSource)
	require.Equal(t, model.DataSourceOpenline, locationWithoutAddressDtls.Source)
	require.Equal(t, "", *locationWithoutAddressDtls.Country)
	require.Equal(t, "", *locationWithoutAddressDtls.Region)
	require.Equal(t, "", *locationWithoutAddressDtls.Locality)
	require.Equal(t, "", *locationWithoutAddressDtls.Address)
	require.Equal(t, "", *locationWithoutAddressDtls.Address2)
	require.Equal(t, "", *locationWithoutAddressDtls.Zip)
	require.False(t, *locationWithoutAddressDtls.Commercial)
	require.Nil(t, locationWithoutAddressDtls.Latitude)
	require.Nil(t, locationWithoutAddressDtls.Longitude)
}

func TestQueryResolver_Contacts_SortByTitleAscFirstNameAscLastNameDesc(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contact1 := neo4jt.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{
		Prefix:    "MR",
		FirstName: "contact",
		LastName:  "1",
	})
	contact2 := neo4jt.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{
		Prefix:    "DR",
		FirstName: "contact",
		LastName:  "9",
	})
	contact3 := neo4jt.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{
		Prefix:    "",
		FirstName: "contact",
		LastName:  "222",
	})
	contact4 := neo4jt.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{
		Prefix:    "MR",
		FirstName: "other contact",
		LastName:  "A",
	})

	rawResponse, err := c.RawPost(getQuery("contact/get_contacts_with_sorting"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contacts struct {
		Contacts model.ContactsPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contacts)
	require.Nil(t, err)
	require.NotNil(t, contacts.Contacts)
	require.Equal(t, 4, len(contacts.Contacts.Content))
	require.Equal(t, contact3, contacts.Contacts.Content[0].ID)
	require.Equal(t, contact2, contacts.Contacts.Content[1].ID)
	require.Equal(t, contact1, contacts.Contacts.Content[2].ID)
	require.Equal(t, contact4, contacts.Contacts.Content[3].ID)

	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
}

func TestQueryResolver_Contact_BasicFilters_FindContactWithLetterAInName(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contactFoundByFirstName := neo4jt.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{
		Prefix:    "MR",
		Name:      "contact1",
		FirstName: "aa",
		LastName:  "bb",
	})
	contactFoundByLastName := neo4jt.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{
		Prefix:    "MR",
		FirstName: "bb",
		LastName:  "AA",
	})
	contactFilteredOut := neo4jt.CreateContact(ctx, driver, tenantName, neo4jentity.ContactEntity{
		Prefix:    "MR",
		FirstName: "bb",
		LastName:  "BB",
	})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contacts_basic_filters"))
	assertRawResponseSuccess(t, rawResponse, err)

	var contactsStruct struct {
		Contacts model.ContactsPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &contactsStruct)
	require.Nil(t, err)
	require.NotNil(t, contactsStruct.Contacts)

	contacts := contactsStruct.Contacts.Content

	require.Equal(t, 2, len(contacts))
	require.Equal(t, contactFoundByFirstName, contacts[0].ID)
	require.Equal(t, "contact1", *contacts[0].Name)
	require.Equal(t, "aa", *contacts[0].FirstName)
	require.Equal(t, "bb", *contacts[0].LastName)
	require.Equal(t, contactFoundByLastName, contacts[1].ID)
	require.Equal(t, 1, contactsStruct.Contacts.TotalPages)
	require.Equal(t, int64(2), contactsStruct.Contacts.TotalElements)

	// suppress unused warnings
	require.NotNil(t, contactFilteredOut)
}

func TestQueryResolver_Contact_WithTimelineEvents(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jtest.CreateDefaultUser(ctx, driver, tenantName)

	now := time.Now().UTC()
	secAgo1 := now.Add(time.Duration(-1) * time.Second)
	secAgo10 := now.Add(time.Duration(-10) * time.Second)
	secAgo40 := now.Add(time.Duration(-40) * time.Second)
	secAgo50 := now.Add(time.Duration(-50) * time.Second)
	secAgo60 := now.Add(time.Duration(-60) * time.Second)

	// prepare page views
	neo4jt.CreatePageView(ctx, driver, contactId, entity.PageViewEntity{
		StartedAt:      secAgo1,
		EndedAt:        now,
		TrackerName:    "tracker1",
		SessionId:      "session1",
		Application:    "application1",
		PageTitle:      "page1",
		PageUrl:        "http://app-1.ai",
		OrderInSession: 1,
		EngagedTime:    10,
	})
	neo4jt.CreatePageView(ctx, driver, contactId, entity.PageViewEntity{
		StartedAt:      secAgo10,
		EndedAt:        now,
		TrackerName:    "tracker2",
		SessionId:      "session2",
		Application:    "application2",
		PageTitle:      "page2",
		PageUrl:        "http://app-2.ai",
		OrderInSession: 2,
		EngagedTime:    20,
	})
	neo4jt.CreatePageView(ctx, driver, contactId2, entity.PageViewEntity{})

	voiceSession := neo4jtest.CreateInteractionSession(ctx, driver, tenantName, "mySessionIdentifier", "session1", "CALL", "ACTIVE", "VOICE", now, false)

	// prepare meeting
	meetingId := neo4jt.CreateMeeting(ctx, driver, tenantName, "meeting-name", secAgo60)
	neo4jt.MeetingAttendedBy(ctx, driver, meetingId, contactId)

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 1", "application/json", channel, secAgo40)
	interactionEventId2 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text 2", "application/json", channel, secAgo50)
	emailId := neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId, "email1", false, "WORK")
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId, "+1234", false, "WORK")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, phoneNumberId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId2, phoneNumberId, "")
	neo4jt.InteractionSessionAttendedBy(ctx, driver, tenantName, voiceSession, phoneNumberId, "")

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "PageView"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Meeting"))
	require.Equal(t, 6, neo4jtest.GetCountOfNodes(ctx, driver, "TimelineEvent"))

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_timeline_events"),
		client.Var("contactId", contactId),
		client.Var("from", now),
		client.Var("size", 8))
	assertRawResponseSuccess(t, rawResponse, err)

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])

	timelineEvents := contact.(map[string]interface{})["timelineEvents"].([]interface{})
	require.Equal(t, 3, len(timelineEvents))

	timelineEvent1 := timelineEvents[0].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent1["__typename"].(string))
	require.Equal(t, interactionEventId1, timelineEvent1["id"].(string))
	require.NotNil(t, timelineEvent1["createdAt"].(string))
	require.Equal(t, "IE text 1", timelineEvent1["content"].(string))
	require.Equal(t, "application/json", timelineEvent1["contentType"].(string))
	require.Equal(t, "EMAIL", timelineEvent1["channel"].(string))

	timelineEvent2 := timelineEvents[1].(map[string]interface{})
	require.Equal(t, "InteractionEvent", timelineEvent2["__typename"].(string))
	require.Equal(t, interactionEventId2, timelineEvent2["id"].(string))
	require.NotNil(t, timelineEvent2["createdAt"].(string))
	require.Equal(t, "IE text 2", timelineEvent2["content"].(string))
	require.Equal(t, "application/json", timelineEvent2["contentType"].(string))
	require.Equal(t, "EMAIL", timelineEvent2["channel"].(string))

	timelineEvent3 := timelineEvents[2].(map[string]interface{})
	require.Equal(t, "Meeting", timelineEvent3["__typename"].(string))
	require.Equal(t, meetingId, timelineEvent3["id"].(string))
	require.NotNil(t, timelineEvent3["createdAt"].(string))
	require.Equal(t, "meeting-name", timelineEvent3["name"].(string))
}

func TestQueryResolver_Contact_WithTimelineEvents_FilterByType(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	now := time.Now().UTC()
	secAgo1 := now.Add(time.Duration(-1) * time.Second)

	neo4jt.CreatePageView(ctx, driver, contactId, entity.PageViewEntity{
		StartedAt:      secAgo1,
		EndedAt:        now,
		TrackerName:    "tracker1",
		SessionId:      "session1",
		Application:    "application1",
		PageTitle:      "page1",
		PageUrl:        "http://app-1.ai",
		OrderInSession: 1,
		EngagedTime:    10,
	})

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "TimelineEvent"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "PageView"))

	types := []model.TimelineEventType{}
	types = append(types, model.TimelineEventTypePageView)

	rawResponse, err := c.RawPost(getQuery("contact/get_contact_with_timeline_filter_by_timeline_event_type"),
		client.Var("contactId", contactId),
		client.Var("from", now),
		client.Var("types", types))
	assertRawResponseSuccess(t, rawResponse, err)

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])

	timelineEvents := contact.(map[string]interface{})["timelineEvents"].([]interface{})
	require.Equal(t, 0, len(timelineEvents))
	//timelineEvent1 := timelineEvents[0].(map[string]interface{})
	//require.Equal(t, "PageView", timelineEvent1["__typename"].(string))
	//require.Equal(t, actionId1, timelineEvent1["id"].(string))
}

func TestQueryResolver_Contact_WithTimelineEventsTotalCount(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	neo4jtest.CreateDefaultUser(ctx, driver, tenantName)

	now := time.Now().UTC()

	// prepare page views
	neo4jt.CreatePageView(ctx, driver, contactId, entity.PageViewEntity{
		StartedAt: now,
	})
	neo4jt.CreatePageView(ctx, driver, contactId, entity.PageViewEntity{
		StartedAt: now,
	})

	// prepare contact notes
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 1", "text/plain", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 2", "text/plain", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 3", "text/plain", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 4", "text/plain", now)
	neo4jt.CreateNoteForContact(ctx, driver, tenantName, contactId, "contact note 5", "text/plain", now)

	// prepare interaction events
	channel := "EMAIL"
	interactionEventId1 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text", "application/json", channel, now)
	interactionEventId2 := neo4jtest.CreateInteractionEvent(ctx, driver, tenantName, "myExternalId", "IE text", "application/json", channel, now)
	emailId := neo4jt.AddEmailTo(ctx, driver, commonModel.CONTACT, tenantName, contactId, "email1", false, "WORK")
	phoneNumberId := neo4jt.AddPhoneNumberTo(ctx, driver, tenantName, contactId, "+1234", false, "WORK")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId1, emailId, "")
	neo4jt.InteractionEventSentTo(ctx, driver, interactionEventId1, phoneNumberId, "")
	neo4jt.InteractionEventSentBy(ctx, driver, interactionEventId2, emailId, "")

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "User"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "PageView"))
	require.Equal(t, 5, neo4jtest.GetCountOfNodes(ctx, driver, "Note"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Email"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "PhoneNumber"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "InteractionEvent"))
	require.Equal(t, 9, neo4jtest.GetCountOfNodes(ctx, driver, "TimelineEvent"))

	rawResponse := callGraphQL(t, "contact/get_contact_with_timeline_events_total_count", map[string]interface{}{"contactId": contactId})

	contact := rawResponse.Data.(map[string]interface{})["contact"]
	require.Equal(t, contactId, contact.(map[string]interface{})["id"])
	require.Equal(t, float64(2), contact.(map[string]interface{})["timelineEventsTotalCount"].(float64))
}

func TestQueryResolver_Contact_WithOrganizations_ById(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	contactId2 := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	organizationId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "organization1")
	organizationId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "organization2")
	organizationId3 := neo4jt.CreateOrganization(ctx, driver, tenantName, "organization3")
	organizationId0 := neo4jt.CreateOrganization(ctx, driver, tenantName, "organization0")
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId, organizationId1)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId, organizationId2)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId, organizationId3)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId2, organizationId0)

	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 4, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 4, neo4jtest.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 4, neo4jtest.GetCountOfRelationships(ctx, driver, "ROLE_IN"))

	rawResponse := callGraphQL(t, "contact/get_contact_with_organizations_by_id",
		map[string]interface{}{"contactId": contactId, "limit": 2, "page": 1})

	var searchedContact struct {
		Contact model.Contact
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &searchedContact)
	require.Nil(t, err)
	require.Equal(t, contactId, searchedContact.Contact.ID)
	require.Equal(t, 2, searchedContact.Contact.Organizations.TotalPages)
	require.Equal(t, int64(3), searchedContact.Contact.Organizations.TotalElements)

	organizations := searchedContact.Contact.Organizations.Content
	require.Equal(t, 2, len(organizations))
	require.Equal(t, organizationId1, organizations[0].ID)
	require.Equal(t, organizationId2, organizations[1].ID)
}

func TestMutationResolver_ContactAddOrganizationByID(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	orgId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org1")
	orgId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org2")
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId, orgId1)

	contactServiceCallbacks := events_platform.MockContactServiceCallbacks{
		LinkWithOrganization: func(context context.Context, request *contactpb.LinkWithOrganizationGrpcRequest) (*contactpb.ContactIdGrpcResponse, error) {
			require.Equal(t, contactId, request.ContactId)
			require.Equal(t, orgId2, request.OrganizationId)
			return &contactpb.ContactIdGrpcResponse{
				Id: contactId,
			}, nil
		},
	}
	events_platform.SetContactCallbacks(&contactServiceCallbacks)

	rawResponse := callGraphQL(t, "contact/add_organization_to_contact", map[string]interface{}{"contactId": contactId, "organizationId": orgId2})

	var contactStruct struct {
		Contact_AddOrganizationById model.Contact
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)
	require.NotNil(t, contactStruct)
}

func TestMutationResolver_ContactRemoveOrganizationByID(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)
	orgId1 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org1")
	orgId2 := neo4jt.CreateOrganization(ctx, driver, tenantName, "org2")
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId, orgId1)
	neo4jt.LinkContactWithOrganization(ctx, driver, contactId, orgId2)

	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "ROLE_IN"))

	callbacks := events_platform.MockEventCompletionCallbacks{
		NotifyEventProcessed: func(context context.Context, org *eventcompletionpb.NotifyEventProcessedRequest) (*emptypb.Empty, error) {
			return &emptypb.Empty{}, nil
		},
	}
	events_platform.SetEventCompletionServiceCallbacks(&callbacks)

	rawResponse := callGraphQL(t, "contact/remove_organization_from_contact", map[string]interface{}{"contactId": contactId, "organizationId": orgId2})

	var contactStruct struct {
		Contact_RemoveOrganizationById model.Contact
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)
	require.NotNil(t, contactStruct)
	organizations := contactStruct.Contact_RemoveOrganizationById.Organizations.Content
	require.Equal(t, contactId, contactStruct.Contact_RemoveOrganizationById.ID)
	require.NotNil(t, contactStruct.Contact_RemoveOrganizationById.UpdatedAt)
	require.Equal(t, 1, len(organizations))
	require.Equal(t, orgId1, organizations[0].ID)
	require.Equal(t, "org1", organizations[0].Name)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "JobRole"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Organization"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "WORKS_AS"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "ROLE_IN"))
}

func TestMutationResolver_ContactAddNewLocation(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	rawResponse := callGraphQL(t, "contact/add_new_location_to_contact", map[string]interface{}{"contactId": contactId})

	var locationStruct struct {
		Contact_AddNewLocation model.Location
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &locationStruct)
	require.Nil(t, err)
	require.NotNil(t, locationStruct)
	location := locationStruct.Contact_AddNewLocation
	require.NotNil(t, location.ID)
	require.NotNil(t, location.CreatedAt)
	require.NotNil(t, location.UpdatedAt)
	require.Equal(t, constants.AppSourceCustomerOsApi, location.AppSource)
	require.Equal(t, model.DataSourceOpenline, location.Source)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Location"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "ASSOCIATED_WITH"))
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, driver, "LOCATION_BELONGS_TO_TENANT"))
	neo4jtest.AssertNeo4jLabels(ctx, t, driver, []string{"Tenant", "Location", "Location_" + tenantName, "Contact", "Contact_" + tenantName})
}

func TestQueryResolver_Contact_WithSocials(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	contactId := neo4jt.CreateDefaultContact(ctx, driver, tenantName)

	socialId1 := neo4jt.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{
		Url: "url1",
	})
	socialId2 := neo4jt.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{
		Url: "url2",
	})
	neo4jt.LinkSocialWithEntity(ctx, driver, contactId, socialId1)
	neo4jt.LinkSocialWithEntity(ctx, driver, contactId, socialId2)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "Social"))
	require.Equal(t, 2, neo4jtest.GetCountOfRelationships(ctx, driver, "HAS"))

	rawResponse := callGraphQL(t, "contact/get_contact_with_socials", map[string]interface{}{"contactId": contactId})

	var contactStruct struct {
		Contact model.Contact
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &contactStruct)
	require.Nil(t, err)

	contact := contactStruct.Contact
	require.NotNil(t, contact)
	require.Equal(t, 2, len(contact.Socials))

	require.Equal(t, socialId1, contact.Socials[0].ID)
	require.Equal(t, "url1", contact.Socials[0].URL)
	require.NotNil(t, contact.Socials[0].CreatedAt)
	require.NotNil(t, contact.Socials[0].UpdatedAt)
	require.Equal(t, "test", contact.Socials[0].AppSource)

	require.Equal(t, socialId2, contact.Socials[1].ID)
	require.Equal(t, "url2", contact.Socials[1].URL)
	require.NotNil(t, contact.Socials[1].CreatedAt)
	require.NotNil(t, contact.Socials[1].UpdatedAt)
	require.Equal(t, "test", contact.Socials[1].AppSource)
}
