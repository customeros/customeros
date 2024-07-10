package graph

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/test"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/aggregate"
	contactEvents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/event"
	contactModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGraphContactEventHandler_OnContactCreate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contactEventHandler := &ContactEventHandler{
		repositories: testDatabase.Repositories,
	}
	myContactId, _ := uuid.NewUUID()
	contactAggregate := aggregate.NewContactAggregateWithTenantAndID(tenantName, myContactId.String())
	curTime := utils.Now()

	dataFields := contactModels.ContactDataFields{
		FirstName:   "Bob",
		LastName:    "Smith",
		Prefix:      "Mr.",
		Description: "This is a test contact.",
	}
	source :=
		events.Source{Source: "N/A", SourceOfTruth: "N/A", AppSource: "event-processing-platform"}
	event, err := contactEvents.NewContactCreateEvent(contactAggregate, dataFields, source, cmnmod.ExternalSystem{}, curTime, curTime)
	require.Nil(t, err)
	err = contactEventHandler.OnContactCreate(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "Contact"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, testDatabase.Driver, "Contact_"+tenantName), "Incorrect number of Contact_%s nodes in Neo4j", tenantName)
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "CONTACT_BELONGS_TO_TENANT"), "Incorrect number of CONTACT_BELONGS_TO_TENANT relationships in Neo4j")

	dbContactNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, myContactId.String())
	require.Nil(t, err)
	require.NotNil(t, dbContactNode)
	contactProps := utils.GetPropsFromNode(*dbContactNode)

	require.Equal(t, myContactId.String(), utils.GetStringPropOrEmpty(contactProps, "id"))
	require.Equal(t, "Bob", utils.GetStringPropOrEmpty(contactProps, "firstName"))
	require.Equal(t, "Smith", utils.GetStringPropOrEmpty(contactProps, "lastName"))
	require.Equal(t, "Mr.", utils.GetStringPropOrEmpty(contactProps, "prefix"))
	require.Equal(t, "This is a test contact.", utils.GetStringPropOrEmpty(contactProps, "description"))

	require.Equal(t, "event-processing-platform", utils.GetStringPropOrEmpty(contactProps, "appSource"))
}

func TestGraphContactEventHandler_OnLocationLinkToContact(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	contactName := "test_contact_name"
	contactId := neo4jtest.CreateContact(ctx, testDatabase.Driver, tenantName, neo4jentity.ContactEntity{
		Name: contactName,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contact": 1, "Contact_" + tenantName: 1})
	dbNodeAfterContactCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, contactId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterContactCreate)
	propsAfterContactCreate := utils.GetPropsFromNode(*dbNodeAfterContactCreate)
	require.Equal(t, contactId, utils.GetStringPropOrEmpty(propsAfterContactCreate, "id"))

	locationName := "test_location_name"
	locationId := neo4jtest.CreateLocation(ctx, testDatabase.Driver, tenantName, neo4jentity.LocationEntity{
		Name: locationName,
	})

	dbNodeAfterLocationCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Location_"+tenantName, locationId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterLocationCreate)
	propsAfterLocationCreate := utils.GetPropsFromNode(*dbNodeAfterLocationCreate)
	require.Equal(t, locationName, utils.GetStringPropOrEmpty(propsAfterLocationCreate, "name"))

	contactEventHandler := &ContactEventHandler{
		repositories: testDatabase.Repositories,
	}
	orgAggregate := aggregate.NewContactAggregateWithTenantAndID(tenantName, contactId)
	now := utils.Now()
	event, err := contactEvents.NewContactLinkLocationEvent(orgAggregate, locationId, now)
	require.Nil(t, err)
	err = contactEventHandler.OnLocationLinkToContact(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "ASSOCIATED_WITH"), "Incorrect number of ASSOCIATED_WITH relationships in Neo4j")
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contactId, "ASSOCIATED_WITH", locationId)
}

func TestGraphContactEventHandler_OnPhoneNumberLinkToContact(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	contactName := "test_contact_name"
	now := utils.Now()
	contactId := neo4jtest.CreateContact(ctx, testDatabase.Driver, tenantName, neo4jentity.ContactEntity{
		Name:      contactName,
		UpdatedAt: now,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contact": 1, "Contact_" + tenantName: 1})
	dbNodeAfterContactCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, contactId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterContactCreate)
	propsAfterContactCreate := utils.GetPropsFromNode(*dbNodeAfterContactCreate)
	require.Equal(t, contactId, utils.GetStringPropOrEmpty(propsAfterContactCreate, "id"))

	validated := false
	e164 := "+0123456789"
	phoneNumberId := neo4jtest.CreatePhoneNumber(ctx, testDatabase.Driver, tenantName, neo4jentity.PhoneNumberEntity{
		E164:           e164,
		Validated:      &validated,
		RawPhoneNumber: e164,
		Source:         constants.SourceOpenline,
		SourceOfTruth:  constants.SourceOpenline,
		AppSource:      constants.SourceOpenline,
	})

	dbNodeAfterPhoneNumberCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "PhoneNumber_"+tenantName, phoneNumberId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterPhoneNumberCreate)
	propsAfterPhoneNumberCreate := utils.GetPropsFromNode(*dbNodeAfterPhoneNumberCreate)
	require.Equal(t, false, utils.GetBoolPropOrFalse(propsAfterPhoneNumberCreate, "validated"))
	creationTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	require.Equal(t, &creationTime, utils.GetTimePropOrNil(propsAfterPhoneNumberCreate, "updatedAt"))

	contactEventHandler := &ContactEventHandler{
		repositories: testDatabase.Repositories,
	}
	contactAggregate := aggregate.NewContactAggregateWithTenantAndID(tenantName, contactId)
	phoneNumberLabel := "phoneNumberLabel"
	updateTime := utils.Now()
	event, err := contactEvents.NewContactLinkPhoneNumberEvent(contactAggregate, phoneNumberId, phoneNumberLabel, true, updateTime)
	require.Nil(t, err)
	err = contactEventHandler.OnPhoneNumberLinkToContact(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contactId, "HAS", phoneNumberId)

	dbContactNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, contactId)
	require.Nil(t, err)
	require.NotNil(t, dbContactNode)
	contactProps := utils.GetPropsFromNode(*dbContactNode)
	require.Less(t, now, *utils.GetTimePropOrNil(contactProps, "updatedAt"))
}

func TestGraphContactEventHandler_OnEmailLinkToContactLinkToContact(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	contactName := "test_contact_name"
	now := utils.Now()
	contactId := neo4jtest.CreateContact(ctx, testDatabase.Driver, tenantName, neo4jentity.ContactEntity{
		Name:      contactName,
		UpdatedAt: now,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contact": 1, "Contact_" + tenantName: 1})
	dbNodeAfterContactCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, contactId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterContactCreate)
	propsAfterContactCreate := utils.GetPropsFromNode(*dbNodeAfterContactCreate)
	require.Equal(t, contactId, utils.GetStringPropOrEmpty(propsAfterContactCreate, "id"))

	email := "email@website.com"
	emailId := neo4jtest.CreateEmail(ctx, testDatabase.Driver, tenantName, neo4jentity.EmailEntity{
		Email:         email,
		RawEmail:      email,
		Primary:       true,
		Source:        constants.SourceOpenline,
		SourceOfTruth: constants.SourceOpenline,
		AppSource:     constants.SourceOpenline,
	})

	dbNodeAfterEmailCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Email_"+tenantName, emailId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterEmailCreate)
	propsAfterEmailCreate := utils.GetPropsFromNode(*dbNodeAfterEmailCreate)
	require.Equal(t, false, utils.GetBoolPropOrFalse(propsAfterEmailCreate, "primary"))
	creationTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	require.Equal(t, &creationTime, utils.GetTimePropOrNil(propsAfterEmailCreate, "updatedAt"))

	contactEventHandler := &ContactEventHandler{
		repositories: testDatabase.Repositories,
	}
	contactAggregate := aggregate.NewContactAggregateWithTenantAndID(tenantName, contactId)
	emailLabel := "emailLabel"
	updateTime := utils.Now()
	userLinkEmailEvent, err := contactEvents.NewContactLinkEmailEvent(contactAggregate, emailId, emailLabel, true, updateTime)
	require.Nil(t, err)
	err = contactEventHandler.OnEmailLinkToContact(context.Background(), userLinkEmailEvent)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "HAS"), "Incorrect number of HAS relationships in Neo4j")
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contactId, "HAS", emailId)

	dbContactNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, contactId)
	require.Nil(t, err)
	require.NotNil(t, dbContactNode)
	contactProps := utils.GetPropsFromNode(*dbContactNode)
	require.Less(t, now, *utils.GetTimePropOrNil(contactProps, "updatedAt"))
}

func TestGraphContactEventHandler_OnContactLinkToOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	contactName := "test_contact_name"
	now := utils.Now()
	contactId := neo4jtest.CreateContact(ctx, testDatabase.Driver, tenantName, neo4jentity.ContactEntity{
		Name:      contactName,
		UpdatedAt: now,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contact": 1, "Contact_" + tenantName: 1})
	dbNodeAfterContactCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, contactId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterContactCreate)
	propsAfterContactCreate := utils.GetPropsFromNode(*dbNodeAfterContactCreate)
	require.Equal(t, contactId, utils.GetStringPropOrEmpty(propsAfterContactCreate, "id"))

	organizationName := "Test Organization"
	organizationId := neo4jtest.CreateOrganization(ctx, testDatabase.Driver, tenantName, neo4jentity.OrganizationEntity{
		Name: organizationName,
	})

	dbNodeAfterOrganizationCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Organization_"+tenantName, organizationId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterOrganizationCreate)
	propsAfterOrganizationCreate := utils.GetPropsFromNode(*dbNodeAfterOrganizationCreate)
	require.Equal(t, organizationName, *utils.GetStringPropOrNil(propsAfterOrganizationCreate, "name"))

	contactEventHandler := &ContactEventHandler{
		repositories: testDatabase.Repositories,
	}
	contactAggregate := aggregate.NewContactAggregateWithTenantAndID(tenantName, contactId)
	jobTitle := "Test Title"
	jobRoleDescription := "Test Description"
	sourceFields := events.Source{
		Source:        constants.SourceOpenline,
		SourceOfTruth: constants.SourceOpenline,
		AppSource:     constants.SourceOpenline,
	}
	curTime := utils.Now()
	endedAt := curTime.AddDate(2, 0, 0)
	event, err := contactEvents.NewContactLinkWithOrganizationEvent(contactAggregate, organizationId, jobTitle, jobRoleDescription, true, sourceFields, curTime, curTime, &curTime, &endedAt)
	require.Nil(t, err)
	err = contactEventHandler.OnContactLinkToOrganization(context.Background(), event)
	require.Nil(t, err)

	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "WORKS_AS"), "Incorrect number of WORKS_AS relationships in Neo4j")
	require.Equal(t, 1, neo4jtest.GetCountOfRelationships(ctx, testDatabase.Driver, "ROLE_IN"), "Incorrect number of ROLE_IN relationships in Neo4j")
	jobRole, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "JobRole_"+tenantName)
	require.Nil(t, err)
	jobRolesProps := utils.GetPropsFromNode(*jobRole)
	jobRoleId := utils.GetStringPropOrEmpty(jobRolesProps, "id")
	require.NotNil(t, jobRoleId)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, contactId, "WORKS_AS", jobRoleId)
	neo4jtest.AssertRelationship(ctx, t, testDatabase.Driver, jobRoleId, "ROLE_IN", organizationId)
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(jobRolesProps, "source"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(jobRolesProps, "sourceOfTruth"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(jobRolesProps, "appSource"))
	require.Equal(t, jobTitle, utils.GetStringPropOrEmpty(jobRolesProps, "jobTitle"))
	require.Equal(t, jobRoleDescription, utils.GetStringPropOrEmpty(jobRolesProps, "description"))
	require.Equal(t, &curTime, utils.GetTimePropOrNil(jobRolesProps, "startedAt"))
	require.Equal(t, &endedAt, utils.GetTimePropOrNil(jobRolesProps, "endedAt"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(jobRolesProps, "primary"))
	require.Equal(t, &curTime, utils.GetTimePropOrNil(jobRolesProps, "createdAt"))
	test.AssertRecentTime(t, utils.GetTimePropOrNow(jobRolesProps, "updatedAt"))

	dbNodeForContactAfterContactLinkToOrganization, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, contactId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeForContactAfterContactLinkToOrganization)
	propsForContactAfterContactLinkToOrganization := utils.GetPropsFromNode(*dbNodeForContactAfterContactLinkToOrganization)
	test.AssertRecentTime(t, utils.GetTimePropOrNow(propsForContactAfterContactLinkToOrganization, "updatedAt"))
}

func TestGraphContactEventHandler_OnContactUpdate(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)

	contactNameCreate := "Contact Name Create"
	contactFirstNameCreate := "Contact FirstName Create"
	contactLastNameCreate := "Contact LastName Create"
	timezoneCreate := "Europe/Paris"
	profilePhotoUrlCreate := "www.pic.com/create"
	prefixCreate := "Mr."
	descriptionCreate := "Description Create"
	sourceOfTruthCreate := constants.SourceOpenline
	now := utils.Now()
	contactId := neo4jtest.CreateContact(ctx, testDatabase.Driver, tenantName, neo4jentity.ContactEntity{
		Name:            contactNameCreate,
		FirstName:       contactFirstNameCreate,
		LastName:        contactLastNameCreate,
		Timezone:        timezoneCreate,
		ProfilePhotoUrl: profilePhotoUrlCreate,
		Prefix:          prefixCreate,
		Description:     descriptionCreate,
		UpdatedAt:       now,
		SourceOfTruth:   neo4jentity.DataSource(sourceOfTruthCreate),
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{"Contact": 1, "Contact_" + tenantName: 1})
	dbNodeAfterContactCreate, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Contact_"+tenantName, contactId)
	require.Nil(t, err)
	require.NotNil(t, dbNodeAfterContactCreate)
	propsAfterContactCreate := utils.GetPropsFromNode(*dbNodeAfterContactCreate)
	require.Equal(t, contactId, utils.GetStringPropOrEmpty(propsAfterContactCreate, "id"))

	contactEventHandler := &ContactEventHandler{
		repositories: testDatabase.Repositories,
	}
	contactAggregate := aggregate.NewContactAggregateWithTenantAndID(tenantName, contactId)
	source := constants.SourceOpenline
	contactNameUpdate := "Contact Name Update"
	contactFirstNameUpdate := "Contact FirstName Update"
	contactLastNameUpdate := "Contact LastName Update"
	timezoneUpdate := "Europe/Bucharest"
	profilePhotoUrlUpdate := "www.pic.com/update"
	prefixUpdate := "Mrs."
	descriptionUpdate := "Description Update"
	dataFields := contactModels.ContactDataFields{
		Name:            contactNameUpdate,
		FirstName:       contactFirstNameUpdate,
		LastName:        contactLastNameUpdate,
		Timezone:        timezoneUpdate,
		ProfilePhotoUrl: profilePhotoUrlUpdate,
		Prefix:          prefixUpdate,
		Description:     descriptionUpdate,
	}
	curTime := utils.Now()
	event, err := contactEvents.NewContactUpdateEvent(contactAggregate, source, dataFields, cmnmod.ExternalSystem{}, curTime,
		[]string{contactEvents.FieldMaskName,
			contactEvents.FieldMaskFirstName,
			contactEvents.FieldMaskLastName,
			contactEvents.FieldMaskTimezone,
			contactEvents.FieldMaskProfilePhotoUrl,
			contactEvents.FieldMaskPrefix,
			contactEvents.FieldMaskDescription})
	require.Nil(t, err)
	err = contactEventHandler.OnContactUpdate(context.Background(), event)
	require.Nil(t, err)

	contact, err := neo4jtest.GetFirstNodeByLabel(ctx, testDatabase.Driver, "Contact_"+tenantName)
	require.Nil(t, err)
	contactProps := utils.GetPropsFromNode(*contact)
	contactId = utils.GetStringPropOrEmpty(contactProps, "id")
	require.NotNil(t, contactId)
	require.Equal(t, contactNameUpdate, utils.GetStringPropOrEmpty(contactProps, "name"))
	require.Equal(t, contactFirstNameUpdate, utils.GetStringPropOrEmpty(contactProps, "firstName"))
	require.Equal(t, contactLastNameUpdate, utils.GetStringPropOrEmpty(contactProps, "lastName"))
	require.Equal(t, timezoneUpdate, utils.GetStringPropOrEmpty(contactProps, "timezone"))
	require.Equal(t, profilePhotoUrlUpdate, utils.GetStringPropOrEmpty(contactProps, "profilePhotoUrl"))
	require.Equal(t, prefixUpdate, utils.GetStringPropOrEmpty(contactProps, "prefix"))
	require.Equal(t, descriptionUpdate, utils.GetStringPropOrEmpty(contactProps, "description"))
	require.Less(t, now, utils.GetTimePropOrNow(contactProps, "updatedAt"))
	require.Equal(t, constants.SourceOpenline, utils.GetStringPropOrEmpty(contactProps, "sourceOfTruth"))
	require.Equal(t, true, utils.GetBoolPropOrFalse(contactProps, "syncedWithEventStore"))
}

func TestGraphContactEventHandler_OnSocialAddedToContactV1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx, testDatabase)(t)

	socialId := uuid.New().String()
	socialUrl := "https://www.facebook.com/contact"
	now := utils.Now()
	neo4jtest.CreateTenant(ctx, testDatabase.Driver, tenantName)
	contactId := neo4jtest.CreateContact(ctx, testDatabase.Driver, tenantName, neo4jentity.ContactEntity{})
	neo4jtest.CreateSocial(ctx, testDatabase.Driver, tenantName, neo4jentity.SocialEntity{Url: socialUrl})
	contactEventHandler := &ContactEventHandler{
		repositories: testDatabase.Repositories,
	}
	contactAggregate := aggregate.NewContactAggregateWithTenantAndID(tenantName, contactId)

	event, err := contactEvents.NewContactAddSocialEvent(contactAggregate,
		socialId, socialUrl, "alias", "ext1", 100,
		events.Source{
			Source:        constants.SourceOpenline,
			SourceOfTruth: constants.SourceOpenline,
			AppSource:     "event-processing-platform",
		}, now)
	require.Nil(t, err)

	// EXECUTE
	err = contactEventHandler.OnSocialAddedToContactV1(context.Background(), event)
	require.Nil(t, err)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, testDatabase.Driver, map[string]int{neo4jutil.NodeLabelContact: 1, neo4jutil.NodeLabelWithTenant(neo4jutil.NodeLabelContact, tenantName): 1, neo4jutil.NodeLabelSocial: 2, neo4jutil.NodeLabelSocial + "_" + tenantName: 2})

	dbNode, err := neo4jtest.GetNodeById(ctx, testDatabase.Driver, "Social_"+tenantName, socialId)
	require.Nil(t, err)
	require.NotNil(t, dbNode)

	social := neo4jmapper.MapDbNodeToSocialEntity(dbNode)
	require.Equal(t, socialId, social.Id)
	require.Equal(t, socialUrl, social.Url)
	require.Equal(t, "alias", social.Alias)
	require.Equal(t, "ext1", social.ExternalId)
	require.Equal(t, int64(100), social.FollowersCount)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), social.Source)
	require.Equal(t, neo4jentity.DataSource(constants.SourceOpenline), social.SourceOfTruth)
	require.Equal(t, "event-processing-platform", social.AppSource)
	require.Equal(t, now, social.CreatedAt)
	test.AssertRecentTime(t, social.UpdatedAt)
}
